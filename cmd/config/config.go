package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	FloodControl struct {
		TimeInterval   time.Duration `yaml:"time-interval"`
		CallLimitCount uint          `yaml:"call-limit-count"`
	} `yaml:"flood-control"`
	Redis struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Database int    `yaml:"database"`
	} `yaml:"redis"`
}

func New(filename string) (Config, error) {
	var config Config

	data, err := os.ReadFile(filename)
	if err != nil {
		return config, fmt.Errorf("failed to read file: %w", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return config, nil
}

// type FloodControlConfig struct {
// 	TimeInterval   time.Duration `yaml:"time-interval"`
// 	CallLimitCount uint          `yaml:"call-limit-count"`
// }

// type Config struct {
// 	FloodControlConfig FloodControlConfig
// 	RedisConfig        redisClient.Config
// }

// func NewFloodControlConfig(filePath string) (FloodControlConfig, error) {
// 	var fcCfg FloodControlConfig
// 	data, err := os.ReadFile(filePath)
// 	if err != nil {
// 		return fcCfg, fmt.Errorf("failed to read config file: %v", err)
// 	}

// 	err = yaml.Unmarshal(data, &fcCfg)
// 	if err != nil {
// 		return fcCfg, fmt.Errorf("failed to unmarshal config data: %v", err)
// 	}

// 	return fcCfg, nil
// }

// func NewConfig(filePath string) *Config {
// 	fcCfg, err := NewFloodControlConfig(filePath)
// 	if err != nil {
// 		log.Fatalf("Can't get FloodControlConfig: %v", err)
// 	}

// 	rcCfg, err := redisClient.NewRedisConfig(filePath)
// 	if err != nil {
// 		log.Fatalf("Can't get FloodRedisConfig: %v", err)
// 	}

// 	return &Config{
// 		FloodControlConfig: fcCfg,
// 		RedisConfig:        rcCfg,
// 	}
// }
