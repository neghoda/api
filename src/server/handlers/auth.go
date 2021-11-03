package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/neghoda/api/src/models"
	"github.com/neghoda/api/src/service"
)

type AuthHandler struct {
	service *service.Service
}

func NewAuthHandler(s *service.Service) *AuthHandler {
	return &AuthHandler{
		service: s,
	}
}

// swagger:operation POST /sign-up auth sign_up
//   registed new user
// ---
// parameters:
// - name: sign_up_request
//   in: body
//   required: true
//   schema:
//     $ref: '#/definitions/SignUpRequest'
// responses:
//   '400':
//     description: Bad Request
//     schema:
//       "$ref": "#/definitions/ValidationErr"
//   '500':
//     description: Internal Server Error
//     schema:
//       "$ref": "#/definitions/CommonError"
func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	req := &models.SignUpRequest{}

	err := UnmarshalRequest(r, req)
	if err != nil {
		SendEmptyResponse(w, http.StatusBadRequest)
		return
	}

	if !ValidatePassword(req.Password) || !ValidateEmail(strings.ToLower(req.Email)) {
		SendEmptyResponse(w, http.StatusBadRequest)
		return
	}

	err = h.service.SignUp(r.Context(), req.Email, req.Password)
	if errors.Is(err, models.ErrAlreadyExist) {
		SendEmptyResponse(w, http.StatusBadRequest)
		return
	}
	if err != nil {
		SendEmptyResponse(w, http.StatusInternalServerError)
		return
	}

	SendResponse(w, http.StatusCreated, nil)
}

// swagger:operation POST /login auth login
//   create a session and obtain tokens pair
// ---
// parameters:
// - name: login_request
//   in: body
//   required: true
//   schema:
//     $ref: '#/definitions/LoginRequest'
// responses:
//   '200':
//     description: Fetched
//     schema:
//       "$ref": "#/definitions/LoginResponse"
//   '400':
//     description: Bad Request
//     schema:
//       "$ref": "#/definitions/ValidationErr"
//   '500':
//     description: Internal Server Error
//     schema:
//       "$ref": "#/definitions/CommonError"
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	req := &models.LoginRequest{}

	err := UnmarshalRequest(r, req)
	if err != nil {
		SendEmptyResponse(w, http.StatusBadRequest)
		return
	}

	if !ValidatePassword(req.Password) || !ValidateEmail(strings.ToLower(req.Email)) {
		SendEmptyResponse(w, http.StatusBadRequest)
		return
	}

	res, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		SendHTTPError(w, err)
		return
	}

	SendResponse(w, http.StatusOK, models.LoginResponse{
		TokenPair: res,
	})
}

// swagger:operation DELETE /logout auth logout
//   deactivate user session, move access token to the black list
// ---
// responses:
//   '204':
//     description: Successfully logged out
//   '400':
//     description: Bad Request
//     schema:
//       "$ref": "#/definitions/ValidationErr"
//   '500':
//     description: Internal Server Error
//     schema:
//       "$ref": "#/definitions/CommonError"
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get(AccessTokenHeader)

	jwtAccessToken, err := ParseAuthorizationHeader(accessToken, BearerSchema)
	if err != nil {
		SendEmptyResponse(w, http.StatusBadRequest)

		return
	}

	err = h.service.Logout(r.Context(), jwtAccessToken)
	if err != nil {
		SendHTTPError(w, err)

		return
	}

	SendResponse(w, http.StatusNoContent, nil)
}

// swagger:operation POST /token auth token
//   refresh access token if previous tokens pair was valid
// ---
// parameters:
// - name: TokenPair
//   in: body
//   required: true
//   schema:
//     $ref: '#/definitions/TokenPair'
// responses:
//   '201':
//     description: Created
//     schema:
//       "$ref": "#/definitions/TokenPair"
//   '400':
//     description: Bad Request
//     schema:
//       "$ref": "#/definitions/ValidationErr"
//   '500':
//     description: Internal Server Error
//     schema:
//       "$ref": "#/definitions/CommonError"
func (h *AuthHandler) TokenRefresh(w http.ResponseWriter, r *http.Request) {
	req := &models.TokenPair{}

	err := UnmarshalRequest(r, req)
	if err != nil {
		SendEmptyResponse(w, http.StatusBadRequest)

		return
	}

	res, err := h.service.RefreshToken(r.Context(), req)
	if err != nil {
		SendHTTPError(w, err)

		return
	}

	SendResponse(w, http.StatusCreated, res)
}
