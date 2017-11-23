package cache

import (
	"time"

	"github.com/HariniGB/login-provider/common"
	"github.com/bluele/gcache"
)

type EvictableMap struct {
	cache   gcache.Cache
	timeout time.Duration
}

func NewEvictableMap(size int, timeout time.Duration) *EvictableMap {
	gc := gcache.New(size).
		LRU().
		Build()

	return &EvictableMap{cache: gc, timeout: timeout}
}

func (e *EvictableMap) Insert(key string, value common.UserInfo) error {
	return e.cache.SetWithExpire(key, value, e.timeout)
}

func (e *EvictableMap) Get(key string) *common.UserInfo {
	uRaw, err := e.cache.GetIFPresent(key)
	if err != nil {
		return nil
	}

	if user, ok := uRaw.(common.UserInfo); ok {
		return &user
	}

	return nil
}

func (e *EvictableMap) Delete(key string) bool {
	return e.cache.Remove(key)
}
