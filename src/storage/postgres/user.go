package postgres

import (
	"context"

	"github.com/neghoda/api/src/models"
)

type UserRepo struct {
	*Connector
}

func (c *Connector) NewUserRepo() *UserRepo {
	return &UserRepo{c}
}

func (ur *UserRepo) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User

	err := ur.QueryContext(ctx).
		Model(&user).
		Where("?TableAlias.email = ?", email).
		First()

	return user, toServiceError(err)
}

func (ur *UserRepo) CreateUser(ctx context.Context, user *models.User) error {
	_, err := ur.QueryContext(ctx).Model(user).Insert()

	return toServiceError(err)
}

func (ur *UserRepo) EmailTaken(ctx context.Context, email string) (bool, error) {
	exist, err := ur.QueryContext(ctx).
		Model((*models.User)(nil)).
		Where("?TableAlias.email = ?", email).
		Exists()

	return exist, toServiceError(err)
}
