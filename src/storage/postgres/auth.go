package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/neghoda/api/src/models"
)

type AuthRepo struct {
	*Connector
}

func (c *Connector) NewAuthRepo() *AuthRepo {
	return &AuthRepo{c}
}

func (ar *AuthRepo) CreateSession(ctx context.Context, session *models.UserSession) error {
	_, err := ar.QueryContext(ctx).Model(session).Insert()

	return err
}

func (ar *AuthRepo) DisableSessionByID(ctx context.Context, sessionID uuid.UUID) error {
	expiredAt := time.Now().UTC()
	session := models.UserSession{
		ID:        sessionID,
		ExpiredAt: &expiredAt,
	}
	_, err := ar.QueryContext(ctx).Model(&session).WherePK().UpdateNotNull()

	return err
}

func (ar *AuthRepo) GetSessionByTokenID(ctx context.Context, tokenID uuid.UUID) (*models.UserSession, error) {
	session := &models.UserSession{}
	err := ar.QueryContext(ctx).
		Model(session).
		Where("token_id = ?", tokenID).
		Limit(1).
		Select()

	return session, err
}

func (ar *AuthRepo) UpdateSession(ctx context.Context, userSession *models.UserSession) error {
	_, err := ar.QueryContext(ctx).
		Model(userSession).
		WherePK().
		Update()

	return err
}
