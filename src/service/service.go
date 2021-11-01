package service

import (
	"context"
	"sync"

	"github.com/erp/api/src/config"
	"github.com/erp/api/src/models"
	"github.com/erp/api/src/storage/redis"
	"github.com/google/uuid"
)

var (
	service *Service
	once    sync.Once
)

type Service struct {
	cfg      *config.Config
	redis    *redis.Client
	authRepo AuthRepo
	userRepo UserRepo
}

func New(
	cfg *config.Config,
	rds *redis.Client,
	aur AuthRepo,
	pr UserRepo,
) *Service {
	once.Do(func() {
		service = &Service{
			redis:    rds,
			cfg:      cfg,
			authRepo: aur,
			userRepo: pr,
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

type UserRepo interface {
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	EmailTaken(ctx context.Context, email string) (bool, error)
}
