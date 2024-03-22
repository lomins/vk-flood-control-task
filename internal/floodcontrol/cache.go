package floodcontrol

import (
	"context"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key, value string, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}
