package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/neghoda/api/src/models"
)

func (q DBQuery) CreateSession(ctx context.Context, session *models.UserSession) error {
	_, err := q.Model(session).Insert()

	return err
}

func (q DBQuery) DisableSessionByID(sessionID uuid.UUID) error {
	expiredAt := time.Now().UTC()
	session := models.UserSession{
		ID:        sessionID,
		ExpiredAt: &expiredAt,
	}
	_, err := q.Model(&session).WherePK().UpdateNotNull()

	return err
}

func (q DBQuery) GetSessionByTokenID(tokenID uuid.UUID) (*models.UserSession, error) {
	session := &models.UserSession{}
	err := q.Model(session).
		Where("token_id = ?", tokenID).
		Limit(1).
		Select()

	return session, err
}

func (q DBQuery) UpdateSession(userSession *models.UserSession) error {
	_, err := q.Model(userSession).WherePK().Update()

	return err
}
