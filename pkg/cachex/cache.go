package cachex

import (
	"context"
	"time"

	"github.com/patrickmn/go-cache"
)

type Cacher interface {
	Set(ctx context.Context, ns, key, value string, expiration ...time.Duration) error
}

var defaultDelimiter = ":"

type options struct {
	Delimiter string
}

type Option func(*options)
