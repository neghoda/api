package context

import (
	"context"

	"github.com/google/uuid"
)

const (
	userIDKey   = "device_id"
	localeKey   = "locale"
	userIPKey   = "user_ip"
	userRoleKey = "user_role"
)

// GetUserID - returns User UID from context.
func GetUserID(ctx context.Context) uuid.UUID {
	value, _ := ctx.Value(userIDKey).(uuid.UUID)

	return value
}

// WithUserID - add User UID value to context.
func WithUserID(ctx context.Context, value uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, value) //nolint
}

// GetLocale - returns Locale from context.
func GetLocale(ctx context.Context) string {
	value, _ := ctx.Value(localeKey).(string)

	return value
}

// WithLocale - add Locale value to context.
func WithLocale(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, localeKey, value) //nolint
}

// GetUserIP - returns User IP from context.
func GetUserIP(ctx context.Context) string {
	value, _ := ctx.Value(userIPKey).(string)

	return value
}

// WithUserIP - add User IP value to context.
func WithUserIP(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, userIPKey, value) //nolint
}
