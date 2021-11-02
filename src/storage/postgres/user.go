package postgres

import (
	"github.com/erp/api/src/models"
)

func (q DBQuery) GetUserByEmail(email string) (models.User, error) {
	var user models.User

	err := q.Model(&user).
		Where("?TableAlias.email = ?", email).
		First()

	return user, toServiceError(err)
}

func (q DBQuery) CreateUser(user *models.User) error {
	_, err := q.Model(user).Insert()

	return toServiceError(err)
}

func (q DBQuery) EmailTaken(email string) (bool, error) {
	exist, err := q.Model((*models.User)(nil)).
		Where("?TableAlias.email = ?", email).
		Exists()

	return exist, toServiceError(err)
}
