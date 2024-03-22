package floodcontrol

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"task/cmd/config"
	"task/internal/storage"
	"time"
)

var (
	ErrCheckCallLimitExceed = errors.New("call limit exceeded")
)

type FloodControl struct {
	TimeInterval   time.Duration
	CallLimitCount uint
	Cache          Cache
}

func New(Cache Cache, config config.Config) *FloodControl {
	log.Println("TimeInterval = ", config.FloodControl.TimeInterval)
	log.Println("CallLimitCount = ", config.FloodControl.CallLimitCount)
	return &FloodControl{
		TimeInterval:   config.FloodControl.TimeInterval,
		CallLimitCount: config.FloodControl.CallLimitCount,
		Cache:          Cache,
	}
}

func (f *FloodControl) Check(ctx context.Context, userID int64) (bool, error) {
	if f.CallLimitCount == 0 {
		return false, ErrCheckCallLimitExceed
	}

	floodControlData, err := f.getDataFromCache(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrRedisKeyNotFound) {
			err = f.setDataInCache(ctx, userID, time.Now(), 1)
			if err != nil {
				return false, fmt.Errorf("failed to check: %w", err)
			}

			return true, nil
		}
		return false, fmt.Errorf("failed to check: %w", err)
	}

	if f.intervalHasExpired(floodControlData.LastCallTime) {
		err = f.setDataInCache(ctx, userID, time.Now(), 1)
		if err != nil {
			return false, fmt.Errorf("failed to check: %w", err)
		}
		return true, nil
	}

	callCount := floodControlData.CallCount + 1

	if f.callLimitHasExceeded(callCount) {
		return false, ErrCheckCallLimitExceed
	}

	err = f.setDataInCache(ctx, userID, floodControlData.LastCallTime, callCount)
	if err != nil {
		return false, fmt.Errorf("failed to check: %w", err)
	}

	return true, nil
}

func (f *FloodControl) getDataFromCache(ctx context.Context, userId int64) (*Payload, error) {
	payloadStr, err := f.Cache.Get(ctx, strconv.Itoa(int(userId)))
	if err != nil {
		return nil, fmt.Errorf("failed to get cached flood control data: %w", err)
	}
	var payload Payload
	err = json.Unmarshal([]byte(payloadStr), &payload)
	if err != nil {
		return nil, fmt.Errorf("failed to get cached flood control data: %w", err)
	}
	return &payload, nil
}

func (f *FloodControl) setDataInCache(ctx context.Context, userId int64, lastCallTime time.Time, callCount uint) error {
	payloadBytes, err := json.Marshal(&Payload{
		LastCallTime: lastCallTime,
		CallCount:    callCount,
	})
	if err != nil {
		return fmt.Errorf("failed to cache flood control data: %w", err)
	}
	err = f.Cache.Set(ctx, strconv.Itoa(int(userId)), string(payloadBytes), f.TimeInterval)
	if err != nil {
		return fmt.Errorf("failed to cache flood control data: %w", err)
	}
	return nil
}

func (f *FloodControl) intervalHasExpired(lastCallAt time.Time) bool {
	return time.Now().After(lastCallAt.Add(f.TimeInterval))
}

func (f *FloodControl) callLimitHasExceeded(callCount uint) bool {
	return callCount > f.CallLimitCount
}
