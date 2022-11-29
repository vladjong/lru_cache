package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vladjong/lru_cache/cache"
)

const (
	parallelFactor = 100_000
	size           = 10
)

func TestAsync(t *testing.T) {
	t.Parallel()
	lruCache := cache.New(size)
	t.Run("correctly stored value", func(t *testing.T) {
		key := "k"
		value := 10
		err := lruCache.Add(key, value)
		if err != nil {
			t.Errorf("Expected Set %v", err)
		}
		storedValue, err := lruCache.Get(key)
		assert.NoError(t, err)
		assert.Equal(t, value, storedValue)
		cap := lruCache.Cap()
		assert.Equal(t, 1, cap)
		err = lruCache.Remove(key)
		assert.NoError(t, err)
		err = lruCache.Add(key, value)
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

func Test_Add(t *testing.T) {
	cache := cache.New(5)
	cache.Add("test", 1)
	cache.Add("test", 2)
	cache.Add("test1", 2)
	cache.Add("test2", 2)

	val, err := cache.Get("test")
	assert.Nil(t, err)

	assert.Equal(t, 2, val)
}

func Test_Cap(t *testing.T) {
	cache := cache.New(5)
	cache.Add("test", 1)
	cache.Add("test", 1)
	cache.Add("test", 1)
	cache.Add("test", 1)
	cache.Add("test", 1)

	assert.Equal(t, 1, cache.Cap())
}
