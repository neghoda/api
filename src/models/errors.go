package models

import "errors"

var (
	ErrInvalidData  = errors.New("invalid data")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrNotFound     = errors.New("not found")
	ErrAlreadyExist = errors.New("already exist")
	ErrConflict     = errors.New("conflict")

	ErrSessionNotFound  = errors.New("session by token ID was not found")
	ErrTokensMismatched = errors.New("access and refresh tokens are mismatched")
	ErrSessionExpired   = errors.New("session expired")

	ErrTokenInvalid       = errors.New("token is invalid")
	ErrTokenInBlackList   = errors.New("token in black list")
	ErrTokenClaimsInvalid = errors.New("token claims are invalid")
)

type InternalError string

func (e InternalError) Error() string {
	return string(e)
}

type BadRequest struct {
	Msg    string `json:"-"`
	Errors []FieldError
}

func (e BadRequest) Error() string {
	return e.Msg
}

type FieldError struct {
	Field string
	Code  string
}
