package main

import (
	"context"
	"log"
	"math/rand"
	"task/cmd/config"
	"task/internal/floodcontrol"
	"task/internal/storage/redisClient"
	"time"
)

func main() {
	cfg, err := config.New("./configs/config.yml")
	if err != nil {
		log.Fatalf("err with cfg: %v", err)
	}

	redisStorage := redisClient.NewRedisStorage(cfg)

	controller := floodcontrol.New(redisStorage, cfg)

	runFloodControl(controller)
}

// FloodControl интерфейс, который нужно реализовать.
// Рекомендуем создать директорию-пакет, в которой будет находиться реализация.
type FloodControl interface {
	// Check возвращает false если достигнут лимит максимально разрешенного
	// кол-ва запросов согласно заданным правилам флуд контроля.
	Check(ctx context.Context, userID int64) (bool, error)
}

func runFloodControl(floodControl FloodControl) {
	log.Println("Simulating flood control...")

	ctx := context.Background()
	for {
		userId := int64(rand.Intn(4))
		is_ok, err := floodControl.Check(ctx, userId)
		if err != nil {
			log.Println("err on check: ", err)
		}
		log.Printf("Received request from user #%d", userId)

		if !is_ok {
			log.Printf("DETECTED FLOODING FROM USER #%d!!!", userId)
		}

		time.Sleep(time.Duration(rand.Intn(int(3))) * time.Second)
	}
}
