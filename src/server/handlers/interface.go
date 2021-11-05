package handlers

import (
	"context"

	"github.com/neghoda/api/src/models"
)

type AuthService interface {
	SignUp(ctx context.Context, sighUpReq *models.SignUpRequest) error
	Login(ctx context.Context, loginReq *models.LoginRequest) (resp models.TokenPair, err error)
	Logout(ctx context.Context, accessToken string) (err error)
	RefreshToken(ctx context.Context, tokenReq *models.TokenPair) (resp models.TokenPair, err error)
}

type FundService interface {
	TickerList() ([]string, error)
	FundByTicker(ctx context.Context, ticker string) (models.Fund, error)
}
