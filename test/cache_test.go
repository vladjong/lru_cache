package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vladjong/lru_cache/cache"
)

const (
	parallelFactor = 100_000
	cappacity      = 10
)

func TestAsync(t *testing.T) {
	t.Parallel()
	lruCache := cache.New(cappacity)
	t.Run("correctly stored value", func(t *testing.T) {
		key := "k"
		value := 10
		val := cache.NewNode(key, value)
		err := lruCache.Add(key, val)
		if err != nil {
			t.Errorf("Expected Set %v", err)
		}
		storedValue, err := lruCache.Get(key)
		assert.NoError(t, err)
		assert.Equal(t, val, storedValue.(*cache.Node).Value)
		cap := lruCache.Cap()
		assert.Equal(t, 10, cap)
		err = lruCache.Remove(key)
		assert.NoError(t, err)
		err = lruCache.Add(key, val)
		if err != nil {
			t.Errorf("Expected Set %v", err)
		}
		lruCache.Clear()
		_, err = lruCache.Get(key)
		assert.Error(t, err, cache.ErrNotFound)
	})
	t.Run("no data races", func(t *testing.T) {
		t.Parallel()
		emulatedLoad(t, lruCache, parallelFactor)
	})
}
