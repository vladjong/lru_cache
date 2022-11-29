package test

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vladjong/lru_cache/cache"
)

const TimeTtl = 100_000_000

func emulatedLoad(t *testing.T, c cache.ICache, parallelFactor int) {
	wg := sync.WaitGroup{}
	for i := 0; i < parallelFactor; i++ {
		key := fmt.Sprintf("%d-key", i)
		value := fmt.Sprintf("%d-val", i)
		wg.Add(1)
		go func(k string) {
			err := c.Add(k, value)
			assert.NoError(t, err)
			wg.Done()
		}(key)
		wg.Add(1)
		go func(k, v string) {
			storedValue, err := c.Get(k)
			if !errors.Is(err, cache.ErrNotFound) {
				assert.Equal(t, v, storedValue)
				assert.NoError(t, err)
			}
			wg.Done()
		}(key, value)
		wg.Add(1)
		go func(k string) {
			err := c.Remove(k)
			if !errors.Is(err, cache.ErrNotFound) {
				assert.NoError(t, err)
			}
			wg.Done()
		}(key)
		wg.Add(1)
		go func(k, v string) {
			c.AddWithTTL(k, value, TimeTtl)
			storedValue, err := c.Get(key)
			if !errors.Is(err, cache.ErrNotFound) {
				assert.Equal(t, v, storedValue)
				assert.NoError(t, err)
			}
			wg.Done()
		}(key, value)
	}
	wg.Wait()
}
