package redisClient

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"task/cmd/config"
	"task/internal/storage"
	"time"

	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database int    `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type RedisStorage struct {
	client *redis.Client
}

func NewRedisConfig(filePath string) (Config, error) {
	var rcCfg Config
	data, err := os.ReadFile(filePath)
	if err != nil {
		return rcCfg, fmt.Errorf("failed to read config file: %v", err)
	}

	err = yaml.Unmarshal(data, &rcCfg)
	if err != nil {
		return rcCfg, fmt.Errorf("failed to unmarshal config data: %v", err)
	}

	rcCfg.User = os.Getenv("REDIS_USER")
	rcCfg.Password = os.Getenv("REDIS_PASSWORD")

	return rcCfg, nil
}

func NewRedisStorage(config config.Config) *RedisStorage {
	redusURL := fmt.Sprintf("redis://%s:%s@%s:%d/%d?protocol=3", config.Redis.User, config.Redis.Password, config.Redis.Host, config.Redis.Port, config.Redis.Database)
	opt, err := redis.ParseURL(redusURL)
	if err != nil {
		log.Fatalf("Error on parse URL: %v", err)
	}

	client := redis.NewClient(opt)

	return &RedisStorage{
		client: client,
	}
}

func (rc *RedisStorage) Set(ctx context.Context, key, value string, expiration time.Duration) error {
	err := rc.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rc *RedisStorage) Get(ctx context.Context, key string) (string, error) {
	val, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", storage.ErrRedisKeyNotFound
		}
		return "", err
	}
	return val, nil
}
