package cache

import "time"

type ICache interface {
	Add(key, value interface{}) error
	Get(key interface{}) (interface{}, error)
	Remove(key interface{}) error
	AddWithTTL(key, value interface{}, ttl time.Duration)
	Cap() int
	Clear()
}
