package memoize

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewMapCache(t *testing.T) {
	cache := NewMapCache(Fn64aKeyer(), WithCleanDur(context.Background(), time.Second))
	if cache == nil {
		t.Errorf("NewMapCache() = %v; want not nil", cache)
	}
}

func TestNewMapCacheWithCleanDur(t *testing.T) {
	ctx := context.Background()
	cache := NewMapCache(Fn64aKeyer(), WithCleanDur(ctx, time.Millisecond))
	assert.Equal(t, time.Millisecond, cache.cleanDur)
	assert.Equal(t, ctx, cache.cleanCtx)
}
