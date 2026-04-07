package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"taoshop/internal/models"
	"taoshop/internal/repository"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService struct {
	users     *repository.UserRepository
	jwtSecret []byte
}

func NewAuthService(users *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		users:     users,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *AuthService) Register(ctx context.Context, username, email, password string) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)//把用户传进来的明文密码 password 进行加密。
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
	}

	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, *models.User, error) {
	user, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil, ErrInvalidCredentials
		}
		return "", nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	signed, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", nil, err
	}

	return signed, user, nil
}

func (s *AuthService) ParseToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return 0, ErrInvalidCredentials
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, ErrInvalidCredentials
	}

	value, ok := claims["sub"].(float64)
	if !ok {
		return 0, ErrInvalidCredentials
	}

	return int64(value), nil
}
