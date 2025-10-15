package redis

import (
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"

	"gitlab.HIDDEN.com/HIDDEN/sms-core/logger"
)

type RedisProvider struct {
	client *redis.Client
}

type LockPayload struct {
	UserID string `json:"userId"`
	Action string `json:"action"`
	Entity string `json:"entity"`
}

type LockKey struct {
	EntityId string `json:"entityId"`
	Entity   string `json:"entity"`
	OrgId    string `json:"orgId"`
}

func NewRedisProvider() RedisProvider {
	logger.Info().Msg("NewRedisProvider!!!")
	redisHostEnv := os.Getenv("REDIS_HOST")
	redisPortEnv := os.Getenv("REDIS_PORT")
	redisPasswordEnv := os.Getenv("REDIS_PASSWORD")

	redisAddr := fmt.Sprintf("%s:%s", redisHostEnv, redisPortEnv)

	redisClient := redis.NewClient(&redis.Options{
		Password: redisPasswordEnv,
		Addr:     redisAddr,
	})

	logger.Info().Msgf("REDIA ADDR: %v", redisAddr)
	logger.Info().Msgf("redisClient %v", redisClient)

	return RedisProvider{
		client: redisClient,
	}
}

func (p *RedisProvider) GetClient() *redis.Client {
	return p.client
}
