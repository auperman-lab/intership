package repository

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type RedisTokenRepo struct {
	Redis *redis.Client
}

func NewTokenRepository(redisClient *redis.Client) *RedisTokenRepo {
	return &RedisTokenRepo{
		Redis: redisClient,
	}
}
func (r *RedisTokenRepo) GetToken(token string) (uint, error) {
	ctx := context.Background()

	idstr, err := r.Redis.Get(ctx, token).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, fmt.Errorf("token not found")
		}
		return 0, fmt.Errorf("error getting token: %w", err)
	}

	id, _ := strconv.ParseUint(idstr, 10, 32)
	if err != nil {
		return 0, err // Return the error if the parsing fails
	}

	return uint(id), nil
}

func (r *RedisTokenRepo) CreateToken(token string, id uint, expiration time.Duration) error {
	ctx := context.Background()

	err := r.Redis.Set(ctx, token, id, expiration).Err()
	if err != nil {
		return fmt.Errorf("error creating token: %w", err)
	}

	return nil
}
