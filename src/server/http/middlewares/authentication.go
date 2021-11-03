package middleware

import (
	"net/http"

	reqContext "github.com/neghoda/api/src/context"
	"github.com/neghoda/api/src/server/handlers"
	"github.com/neghoda/api/src/service"
	log "github.com/sirupsen/logrus"
)

// Auth - authenticate User by JWT token and add to context his ID, role etc.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			authToken = r.Header.Get(handlers.AccessTokenHeader)
		)

		jwtToken, err := handlers.ParseAuthorizationHeader(authToken, handlers.BearerSchema)
		if err != nil {
			log.Error("Token is invalid ", authToken)
			handlers.SendEmptyResponse(w, http.StatusUnauthorized)
			return
		}

		loginSes, err := service.Get().Validate(jwtToken)
		if err != nil {
			log.Error("Token is invalid ", authToken)
			handlers.SendEmptyResponse(w, http.StatusUnauthorized)
			return
		}

		ctx := reqContext.WithUserID(r.Context(), loginSes.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
