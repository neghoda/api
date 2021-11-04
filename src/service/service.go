package service

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/neghoda/api/src/config"
	"github.com/neghoda/api/src/models"
	"github.com/neghoda/api/src/repo/ssga"
)

var (
	service *Service
	once    sync.Once
)

type Service struct {
	cfg      *config.Config
	authRepo AuthRepo
	fundRepo FundRepo
	userRepo UserRepo
	ssgaRepo *ssga.Client
}

func New(
	cfg *config.Config,
	authRepo AuthRepo,
	fundRepo FundRepo,
	userRepo UserRepo,
	ssgaRepo *ssga.Client,
) *Service {
	once.Do(func() {
		service = &Service{
			cfg:      cfg,
			authRepo: authRepo,
			fundRepo: fundRepo,
			userRepo: userRepo,
			ssgaRepo: ssgaRepo,
		}
	})

	return service
}

func Get() *Service {
	return service
}

type AuthRepo interface {
	CreateSession(ctx context.Context, session *models.UserSession) error
	DisableSessionByID(ctx context.Context, sessionID uuid.UUID) error
	GetSessionByTokenID(ctx context.Context, tokenID uuid.UUID) (*models.UserSession, error)
	UpdateSession(ctx context.Context, userSession *models.UserSession) error
}

type FundRepo interface {
	FetchFund(ctx context.Context, ticker string) (models.Fund, error)
	ReplaceFund(ctx context.Context, fund *models.Fund) error
}

type UserRepo interface {
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	EmailTaken(ctx context.Context, email string) (bool, error)
}
