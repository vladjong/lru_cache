package test

import (
	"testing"
	"time"

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
		assert.Error(t, cache.ErrNotFound, err)
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

func Test_Ttl(t *testing.T) {
	c := cache.New(5)
	c.AddWithTTL("test", 1, 10000000000)
	time.Sleep(time.Microsecond * 1)
	val, err := c.Get("test")
	assert.Nil(t, err)
	assert.Equal(t, 1, val)
	time.Sleep(time.Second * 3)
	_, err = c.Get("test")
	assert.Error(t, cache.ErrNotFound, err)
}
