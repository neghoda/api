package service

import (
	"github.com/golang-jwt/jwt"
	"github.com/neghoda/api/src/models"
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

	return &models.UserSession{
		UserID:  claims.UserID,
		TokenID: claims.TokenID,
	}, nil
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
