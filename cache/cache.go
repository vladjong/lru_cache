package cache

type ICache interface {
	Add(key, value interface{}) error
	Get(key interface{}) (interface{}, error)
	Remove(key interface{}) error
	Cap() int
	Clear()
}
