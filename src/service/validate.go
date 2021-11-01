package service

import (
	"github.com/erp/api/src/models"
	"github.com/erp/api/src/storage/redis"
	"github.com/golang-jwt/jwt"
)

// Validate validates access token.
func (s Service) Validate(accessToken string) (*models.UserSession, error) {
	token, err := s.parseJWT(accessToken)
	if err != nil {
		return nil, err
	}

	claims := s.parseClaims(token)
	if claims == nil || !token.Valid {
		return nil, models.ErrTokenInvalid
	}

	if err := s.checkBlackList(claims.TokenID.String()); err != nil {
		return nil, err
	}

	return &models.UserSession{
		UserID:  claims.UserID,
		TokenID: claims.TokenID,
	}, nil
}

func (s Service) checkBlackList(tokenID string) error {
	res, err := s.redis.Get(redis.TokenBlackListKey(tokenID))
	if err != nil {
		return err
	}

	if res != nil {
		return models.ErrTokenInBlackList
	}

	return nil
}

func (s Service) parseJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.AccessTokenSecret), nil
	})

	if token == nil {
		return nil, models.ErrTokenInvalid
	}

	//nolint: errorlint
	switch ve := err.(type) {
	case *jwt.ValidationError:
		if ve.Errors|(jwt.ValidationErrorExpired) != jwt.ValidationErrorExpired {
			return nil, models.ErrTokenInvalid
		}
	case nil:
	default:
		return nil, err
	}

	return token, nil
}
