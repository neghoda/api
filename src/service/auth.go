package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/neghoda/api/src/models"
	"golang.org/x/crypto/bcrypt"
)

const hoursInDay = 24

func (s Service) SignUp(ctx context.Context, signUpReq *models.SignUpRequest) error {
	if !validatePassword(signUpReq.Password) || !validateEmail(strings.ToLower(signUpReq.Email)) {
		return models.ErrInvalidData
	}

	emailTaken, err := s.userRepo.EmailTaken(ctx, signUpReq.Email)
	if err != nil {
		return err
	}

	if emailTaken {
		return models.ErrAlreadyExist
	}

	hashPassword, err := newHash(signUpReq.Password)
	if err != nil {
		return err
	}

	user := &models.User{
		Email:     strings.ToLower(signUpReq.Email),
		Password:  hashPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.userRepo.CreateUser(ctx, user)
}

func (s Service) Login(ctx context.Context, loginReq *models.LoginRequest) (models.TokenPair, error) {
	var tokenPair models.TokenPair

	if !validatePassword(loginReq.Password) || !validateEmail(strings.ToLower(loginReq.Email)) {
		return tokenPair, models.ErrUnauthorized
	}

	user, err := s.userRepo.GetUserByEmail(ctx, loginReq.Email)
	if err != nil {
		return tokenPair, err
	}

	if !compareHashAndPassword(loginReq.Password, user.Password) {
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
func (s Service) RefreshToken(ctx context.Context, tokenReq *models.TokenPair) (models.TokenPair, error) {
	access, err := s.parseJWT(tokenReq.AccessToken)
	if err != nil {
		return models.TokenPair{}, models.ErrTokenInvalid
	}

	accessClaims := s.parseClaims(access)
	if accessClaims == nil {
		return models.TokenPair{}, models.ErrTokenClaimsInvalid
	}

	session, err := s.authRepo.GetSessionByTokenID(ctx, accessClaims.TokenID)
	if err != nil {
		return models.TokenPair{}, models.ErrSessionNotFound
	}

	if session.RefreshToken != tokenReq.RefreshToken {
		return models.TokenPair{}, models.ErrTokensMismatched
	}

	now := time.Now().UTC()
	if now.After(*session.ExpiredAt) {
		return models.TokenPair{}, models.ErrSessionExpired
	}

	claims := models.Claims{
		TokenID:   uuid.New(),
		SessionID: session.ID,
		UserID:    session.UserID,
	}

	accessNew, err := s.GenerateAccess(&claims)
	if err != nil {
		return models.TokenPair{}, err
	}

	session.TokenID = claims.TokenID
	session.UpdatedAt = &now

	err = s.authRepo.UpdateSession(ctx, session)
	if err != nil {
		return models.TokenPair{}, err
	}

	if _, err = s.Revoke(tokenReq.AccessToken); err != nil {
		return models.TokenPair{}, err
	}

	return models.TokenPair{
		AccessToken:  accessNew,
		RefreshToken: tokenReq.RefreshToken,
	}, nil
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

	return claims, nil
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

func validatePassword(password string) bool {
	var (
		hasUpper  = false
		hasLower  = false
		hasNumber = false
	)

	if !(len(password) >= 8) || !(len(password) <= 50) {
		return false
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	return hasUpper && hasLower && hasNumber
}

func validateEmail(email string) bool {
	if m, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, email); !m {
		return false
	}

	return true
}
