package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/xarcher/backend/internal/domain"
)

type JWTService interface {
	GenerateToken(claims *domain.TokenClaims) (string, error)
	ValidateToken(tokenString string) (*domain.TokenClaims, error)
	RevokeToken(token string) error
}

type jwtService struct {
	secretKey     []byte
	revokedTokens map[string]time.Time // In production, use Redis or database
}

func NewJWTService(secretKey string) JWTService {
	return &jwtService{
		secretKey:     []byte(secretKey),
		revokedTokens: make(map[string]time.Time),
	}
}

func (j *jwtService) GenerateToken(claims *domain.TokenClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  claims.UserID,
		"username": claims.Username,
		"iat":      claims.IssuedAt,
		"exp":      claims.ExpiresAt,
	})

	return token.SignedString(j.secretKey)
}

func (j *jwtService) ValidateToken(tokenString string) (*domain.TokenClaims, error) {
	// Check if token is revoked
	if _, revoked := j.revokedTokens[tokenString]; revoked {
		return nil, errors.New("token has been revoked")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return &domain.TokenClaims{
		UserID:    int(claims["user_id"].(float64)),
		Username:  claims["username"].(string),
		IssuedAt:  int64(claims["iat"].(float64)),
		ExpiresAt: int64(claims["exp"].(float64)),
	}, nil
}

func (j *jwtService) RevokeToken(token string) error {
	j.revokedTokens[token] = time.Now()
	return nil
}
