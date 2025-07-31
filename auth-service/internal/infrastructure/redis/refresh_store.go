package redisstore

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RefreshStore struct {
	rdb *redis.Client
}

func NewRefreshStore(addr string, db int) *RefreshStore {
	return &RefreshStore{
		rdb: redis.NewClient(&redis.Options{
			Addr: addr,
			DB:   db,
		}),
	}
}

// сохраняет refresh-токен (по jti) с TTL.
func (s *RefreshStore) SaveRefresh(ctx context.Context, jti string, userID int64, ttl time.Duration) error {
	key := fmt.Sprintf("refresh:%s", jti)
	return s.rdb.Set(ctx, key, userID, ttl).Err()
}

// проверяет, что токен существует.
func (s *RefreshStore) ValidateRefresh(ctx context.Context, jti string) (int64, error) {
	key := fmt.Sprintf("refresh:%s", jti)
	id, err := s.rdb.Get(ctx, key).Int64()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// удаляет refresh-токен (при логауте/ротации).
func (s *RefreshStore) RevokeRefresh(ctx context.Context, jti string) error {
	key := fmt.Sprintf("refresh:%s", jti)
	return s.rdb.Del(ctx, key).Err()
}
