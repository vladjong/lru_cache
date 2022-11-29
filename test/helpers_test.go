package test

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vladjong/lru_cache/cache"
)

func emulatedLoad(t *testing.T, c cache.ICache, parallelFactor int) {
	wg := sync.WaitGroup{}
	for i := 0; i < parallelFactor; i++ {
		key := fmt.Sprintf("%d-key", i)
		value := cache.NewNode(key, i)
		wg.Add(1)
		go func(k string) {
			err := c.Add(k, value)
			assert.NoError(t, err)
			wg.Done()
		}(key)
		wg.Add(1)
		go func(k string, v *cache.Node) {
			storedValue, err := c.Get(k)
			if !errors.Is(err, cache.ErrNotFound) {
				assert.Equal(t, v, storedValue.(*cache.Node).Value)
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
		go func() {
			cap := c.Cap()
			assert.Equal(t, 10, cap)
			wg.Done()
		}()
	}
	wg.Wait()
}
