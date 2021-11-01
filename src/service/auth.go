package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"strings"
	"time"

	"github.com/erp/api/src/models"
	"github.com/erp/api/src/storage/redis"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const hoursInDay = 24

func (s Service) SignUp(ctx context.Context, email, password string) error {
	emailTaken, err := s.userRepo.EmailTaken(ctx, email)
	if err != nil {
		return err
	}

	if emailTaken {
		return models.ErrAlreadyExist
	}

	hashPassword, err := newHash(password)
	if err != nil {
		return err
	}

	user := &models.User{
		Email:     strings.ToLower(email),
		Password:  hashPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.userRepo.CreateUser(ctx, user)
}

func (s Service) Login(ctx context.Context, email, password string) (models.TokenPair, error) {
	var tokenPair models.TokenPair

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return tokenPair, err
	}

	if !compareHashAndPassword(password, user.Password) {
		return tokenPair, models.ErrUnauthorized
	}

	claims := models.Claims{
		SessionID: uuid.New(),
		TokenID:   uuid.New(),
		UserID:    user.ID,
	}

	accessToken, err := s.GenerateAccess(&claims)
	if err != nil {
		return tokenPair, err
	}

	tokenPair.AccessToken = accessToken

	refreshToken, err := s.GenerateRefresh()
	if err != nil {
		return tokenPair, err
	}

	tokenPair.RefreshToken = refreshToken

	err = s.createSession(ctx, &claims, &tokenPair)
	if err != nil {
		return tokenPair, err
	}

	return tokenPair, err
}

func (s Service) Logout(ctx context.Context, accessToken string) (err error) {
	claims, err := s.Revoke(accessToken)
	if err != nil {
		return err
	}

	err = s.authRepo.DisableSessionByID(ctx, claims.SessionID)
	if err != nil {
		return err
	}

	return nil
}

// RefreshToken refreshes access token.
func (s Service) RefreshToken(ctx context.Context, oldTokens *models.TokenPair) (*models.TokenPair, error) {
	access, err := s.parseJWT(oldTokens.AccessToken)
	if err != nil {
		return nil, models.ErrTokenInvalid
	}

	accessClaims := s.parseClaims(access)
	if accessClaims == nil {
		return nil, models.ErrTokenClaimsInvalid
	}

	err = s.checkBlackList(accessClaims.TokenID.String())
	if err != nil {
		return nil, models.ErrTokenInvalid
	}

	session, err := s.authRepo.GetSessionByTokenID(ctx, accessClaims.TokenID)
	if err != nil {
		return nil, models.ErrSessionNotFound
	}

	if session.RefreshToken != oldTokens.RefreshToken {
		return nil, models.ErrTokensMismatched
	}

	now := time.Now().UTC()
	if now.After(*session.ExpiredAt) {
		return nil, models.ErrSessionExpired
	}

	claims := models.Claims{
		TokenID:   uuid.New(),
		SessionID: session.ID,
		UserID:    session.UserID,
	}

	accessNew, err := s.GenerateAccess(&claims)
	if err != nil {
		return nil, err
	}

	session.TokenID = claims.TokenID
	session.UpdatedAt = &now

	err = s.authRepo.UpdateSession(ctx, session)
	if err != nil {
		return nil, err
	}

	if _, err = s.Revoke(oldTokens.AccessToken); err != nil {
		return nil, err
	}

	return &models.TokenPair{AccessToken: accessNew, RefreshToken: oldTokens.RefreshToken}, nil
}

// GenerateAccess generates token with claims.
func (s Service) GenerateAccess(claims *models.Claims) (string, error) {
	claims.StandardClaims.ExpiresAt = time.Now().Unix() + int64(s.cfg.AccessTokenTTL)
	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := tokenWithClaims.SignedString([]byte(s.cfg.AccessTokenSecret))
	if err != nil {
		return "", err
	}

	return token, nil
}

// GenerateRefresh generates refresh token.
func (s Service) GenerateRefresh() (string, error) {
	return generateRandomString(s.cfg.RefreshTokenLen)
}

func (s Service) parseClaims(token *jwt.Token) *models.Claims {
	if claims, ok := token.Claims.(*models.Claims); ok {
		return claims
	}

	return nil
}

// getExpiredAt returns session expiration time till the end of the current day.
func getExpiredAt() time.Time {
	return time.Now().UTC().AddDate(0, 0, 1).Truncate(time.Hour * hoursInDay)
}

func (s Service) createSession(ctx context.Context, claims *models.Claims, tokenPair *models.TokenPair) error {
	now := time.Now().UTC()
	expiredAt := getExpiredAt()
	session := models.UserSession{
		ID:           claims.SessionID,
		UserID:       claims.UserID,
		TokenID:      claims.TokenID,
		RefreshToken: tokenPair.RefreshToken,
		CreatedAt:    &now,
		UpdatedAt:    &now,
		ExpiredAt:    &expiredAt,
	}

	return s.authRepo.CreateSession(ctx, &session)
}

// Revoke revokes access token.
func (s Service) Revoke(accessToken string) (*models.Claims, error) {
	token, err := s.parseJWT(accessToken)
	if err != nil {
		return nil, err
	}

	claims := s.parseClaims(token)
	if claims == nil {
		return nil, models.ErrTokenClaimsInvalid
	}

	err = s.AddToBlacklist(claims.TokenID.String(), claims.TTL())
	if err != nil {
		return nil, err
	}

	return claims, nil
}

// AddToBlacklist adds token to the blacklist.
func (s Service) AddToBlacklist(tokenID string, ttl int64) error {
	if ttl <= 0 {
		ttl = int64(s.cfg.AccessTokenTTL)
	}

	return s.redis.Set(redis.TokenBlackListKey(tokenID), []byte{}, ttl)
}

// generateRandomString returns a URL-safe, base64 encoded securely generated random string.
func generateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}

func newHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func compareHashAndPassword(password, hash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}
	return true
}
