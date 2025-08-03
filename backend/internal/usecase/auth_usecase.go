package usecase

import (
	_ "crypto/rand"
	"errors"
	_ "fmt"
	"time"

	"github.com/xarcher/backend/internal/domain"
	"github.com/xarcher/backend/internal/infrastructure/jwt"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	userRepo   domain.UserRepository
	jwtService jwt.JWTService
	timeout    time.Duration
}

func NewAuthUsecase(userRepo domain.UserRepository, jwtService jwt.JWTService, timeout time.Duration) domain.AuthUsecase {
	return &authUsecase{
		userRepo:   userRepo,
		jwtService: jwtService,
		timeout:    timeout,
	}
}

func (a *authUsecase) Register(req *domain.AuthRequest) (*domain.AuthResponse, error) {
	// Check if user exists
	existingUser, _ := a.userRepo.GetByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &domain.User{
		Username:  req.Username,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	if err := a.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate token
	return a.generateTokenResponse(user)
}

func (a *authUsecase) Login(req *domain.AuthRequest) (*domain.AuthResponse, error) {
	// Get user
	user, err := a.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate token
	return a.generateTokenResponse(user)
}

func (a *authUsecase) ValidateToken(token string) (*domain.TokenClaims, error) {
	return a.jwtService.ValidateToken(token)
}

func (a *authUsecase) RevokeToken(token string) error {
	return a.jwtService.RevokeToken(token)
}

func (a *authUsecase) generateTokenResponse(user *domain.User) (*domain.AuthResponse, error) {
	expiresAt := time.Now().Add(24 * time.Hour)

	claims := &domain.TokenClaims{
		UserID:    user.ID,
		Username:  user.Username,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: expiresAt.Unix(),
	}

	token, err := a.jwtService.GenerateToken(claims)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}
