package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (r *RedisTokenRepo) CreateToken(token string, id primitive.ObjectID, expiration time.Duration) error {
	ctx := context.Background()

	err := r.Redis.Set(ctx, token, id, expiration).Err()
	if err != nil {

		return fmt.Errorf("error creating token: %w", err)
	}

	return nil
}
