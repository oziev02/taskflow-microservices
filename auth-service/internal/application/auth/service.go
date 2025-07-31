package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/oziev02/taskflow-microservices/auth-service/internal/domain/user"
	jwtmgr "github.com/oziev02/taskflow-microservices/auth-service/internal/infrastructure/jwt"
	redisstore "github.com/oziev02/taskflow-microservices/auth-service/internal/infrastructure/redis"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Service struct {
	repo  user.Repository
	jwt   *jwtmgr.Manager
	store *redisstore.RefreshStore
}

func NewService(repo user.Repository, jwt *jwtmgr.Manager, store *redisstore.RefreshStore) *Service {
	return &Service{repo: repo, jwt: jwt, store: store}
}

func (s *Service) Register(ctx context.Context, email, password string) (int64, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("hash password: %w", err)
	}
	return s.repo.Create(email, string(hash))
}

func (s *Service) Login(ctx context.Context, email, password string) (access, refresh string, err error) {
	u, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", "", ErrInvalidCredentials
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return "", "", ErrInvalidCredentials
	}
	at, rt, jti, err := s.jwt.GeneratePair(u.ID, u.Email)
	if err != nil {
		return "", "", err
	}
	// Сохраняем refresh (по jti) в Redis с TTL=refreshTTL
	if err := s.store.SaveRefresh(ctx, jti, u.ID, 720*time.Hour); err != nil {
		return "", "", err
	}
	return at, rt, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	claims, err := s.jwt.Parse(refreshToken)
	if err != nil {
		return "", "", err
	}
	if claims.Subject != "refresh" || claims.ID == "" {
		return "", "", errors.New("invalid refresh token")
	}
	// Проверяем, что refresh существует в Redis
	if _, err := s.store.ValidateRefresh(ctx, claims.ID); err != nil {
		return "", "", errors.New("refresh token not found or expired")
	}
	// Генерируем новую пару
	at, rt, newJTI, err := s.jwt.GeneratePair(claims.UserID, claims.Email)
	if err != nil {
		return "", "", err
	}
	// Ротируем: старый refresh удаляем, новый сохраняем
	_ = s.store.RevokeRefresh(ctx, claims.ID)
	if err := s.store.SaveRefresh(ctx, newJTI, claims.UserID, 720*time.Hour); err != nil {
		return "", "", err
	}
	return at, rt, nil
}
