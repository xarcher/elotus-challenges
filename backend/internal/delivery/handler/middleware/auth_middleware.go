package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/xarcher/backend/internal/domain"
	"github.com/xarcher/backend/pkg/utils"
)

type AuthMiddleware struct {
	authUsecase domain.AuthUsecase
}

func NewAuthMiddleware(authUsecase domain.AuthUsecase) *AuthMiddleware {
	return &AuthMiddleware{
		authUsecase: authUsecase,
	}
}

func (m *AuthMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			utils.RespondError(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		token := tokenParts[1]
		claims, err := m.authUsecase.ValidateToken(token)
		if err != nil {
			utils.RespondError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "username", claims.Username)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
