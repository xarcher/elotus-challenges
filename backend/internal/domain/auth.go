package domain

import "time"

type AuthRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type TokenClaims struct {
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

type RevokedToken struct {
	ID        int       `json:"id" db:"id"`
	Token     string    `json:"token" db:"token"`
	RevokedAt time.Time `json:"revoked_at" db:"revoked_at"`
}

type AuthUsecase interface {
	Register(req *AuthRequest) (*AuthResponse, error)
	Login(req *AuthRequest) (*AuthResponse, error)
	ValidateToken(token string) (*TokenClaims, error)
	RevokeToken(token string) error
}
