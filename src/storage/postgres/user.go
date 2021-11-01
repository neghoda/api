package postgres

import (
	"context"

	"github.com/erp/api/src/models"
)

type UserRepo struct {
	*Postgres
}

func (p *Postgres) NewProfileRepo() *UserRepo {
	return &UserRepo{p}
}

func (p *UserRepo) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User

	err := p.WithContext(ctx).
		Model(&user).
		Where("?TableAlias.email = ?", email).
		First()

	return user, toServiceError(err)
}

func (p *UserRepo) CreateUser(ctx context.Context, user *models.User) error {
	_, err := p.WithContext(ctx).
		Model(user).
		Insert()

	return toServiceError(err)
}

func (p *UserRepo) EmailTaken(ctx context.Context, email string) (bool, error) {
	exist, err := p.WithContext(ctx).
		Model((*models.User)(nil)).
		Where("?TableAlias.email = ?", email).
		Exists()

	return exist, toServiceError(err)
}
