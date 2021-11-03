package handlers

import (
	"context"
	"github.com/neghoda/api/src/models"
)

type AuthService interface {
	Login(ctx context.Context, loginReq models.LoginRequest) (resp models.LoginResponse, err error)
	Logout(ctx context.Context, accessToken string) (err error)
	RefreshToken(ctx context.Context, tokenReq models.TokenPair) (resp models.TokenPair, err error)
}
