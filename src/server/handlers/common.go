package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	log "github.com/sirupsen/logrus"

	"github.com/neghoda/api/src/models"
)

var rxNotAlphanumeric = regexp.MustCompile(`[^a-zA-Z]+`) // to keep Alphanumeric

const (
	contentTypeJSON   = "application/json;charset=utf-8"
	AccessTokenHeader = "Authorization"
	BearerSchema      = "Bearer"
	BasicSchema       = "Basic"
)

// HandleRequest reads all bytes from request and closes it to prevent memory leaks.
func HandleRequest(r *http.Request) (data []byte, err error) {
	data, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return data, r.Body.Close()
}

func SendEmptyResponse(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(statusCode)
}

// SendResponse - common method for encoding and writing any json response.
func SendResponse(w http.ResponseWriter, statusCode int, respBody interface{}) {
	w.Header().Set("Content-Type", contentTypeJSON)

	binRespBody, err := json.Marshal(respBody)
	if err != nil {
		statusCode = http.StatusInternalServerError

		log.Error(err)
	}

	w.WriteHeader(statusCode)
	// nolint
	_, err = w.Write(binRespBody)
	if err != nil {
		log.Error(err)
	}
}

// SendHTTPError sends HTTP error with code.
// nolint: errorlint
func SendHTTPError(w http.ResponseWriter, err error) {
	var (
		statusCode int
		errCode    string
		errMessage string
		errs       []models.FieldError
	)

	switch err {
	case models.ErrUnauthorized,
		models.ErrSessionNotFound,
		models.ErrTokensMismatched,
		models.ErrSessionExpired,
		models.ErrTokenInvalid,
		models.ErrTokenInBlackList,
		models.ErrTokenClaimsInvalid:
		statusCode = http.StatusUnauthorized
		errCode = "unauthorized"
	case models.ErrNotFound:
		statusCode = http.StatusNotFound
		errCode = "not_found"
	case models.ErrForbidden:
		statusCode = http.StatusForbidden
		errCode = "forbidden"
	case models.ErrAlreadyExist,
		models.ErrConflict:
		statusCode = http.StatusConflict
		errCode = "conflict"
	default:
		switch v := err.(type) {
		case models.BadRequest:
			statusCode = http.StatusBadRequest
			errs = v.Errors
			errCode = "bad_request"
			errMessage = v.Msg
		case models.InternalError:
			statusCode = http.StatusInternalServerError
			errCode = "internal"
		default:
			statusCode = http.StatusServiceUnavailable
			errCode = "service_unavailable"
		}
	}

	log.WithError(err).Error("handler error")

	SendResponse(w, statusCode, struct {
		Code    string              `json:"code"`
		Message string              `json:"message,omitempty"`
		Errors  []models.FieldError `json:"validation_errors,omitempty"`
	}{
		Code:    errCode,
		Message: errMessage,
		Errors:  errs,
	})
}

// UnmarshalRequest - unmarshal http request to provided structure.
// nolint: errorlint
func UnmarshalRequest(r *http.Request, body interface{}) error {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return models.InternalError(err.Error())
	}
	defer r.Body.Close()

	if err := json.Unmarshal(reqBody, body); err != nil {
		if e, ok := err.(*json.UnmarshalTypeError); ok {
			return models.BadRequest{
				Msg: err.Error(),
				Errors: []models.FieldError{{
					Field: strings.ToLower(e.Field),
					Code:  "validation_is_" + rxNotAlphanumeric.ReplaceAllString(e.Type.String(), ""),
				}},
			}
		}

		return models.BadRequest{Msg: err.Error()}
	}

	return nil
}

func ParseAuthorizationHeader(header, schema string) (string, error) {
	authSlice := strings.Split(strings.TrimSpace(header), " ")
	if len(authSlice) != 2 || !strings.EqualFold(authSlice[0], schema) {
		return "", errors.New("invalid authorization header") //nolint: goerr113
	}

	return authSlice[1], nil
}

func ValidatePassword(password string) bool {
	var (
		hasUpper  = false
		hasLower  = false
		hasNumber = false
	)

	if !(len(password) >= 8) || !(len(password) <= 50) {
		return false
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	return hasUpper && hasLower && hasNumber
}

func ValidateEmail(email string) bool {
	if m, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, email); !m {
		return false
	}

	return true
}
