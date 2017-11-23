package memory

import (
	"fmt"
	"time"

	"github.com/HariniGB/login-provider/cache"
	"github.com/HariniGB/login-provider/common"
	"github.com/HariniGB/login-provider/storage"
)

type memory struct {
	e *cache.EvictableMap
}

func NewMemoryStorage(size int, timeout time.Duration) storage.Storage {
	e := cache.NewEvictableMap(size, timeout)
	return &memory{e: e}
}

func (m *memory) Get(id string) (*common.UserInfo, error) {
	info := m.e.Get(id)
	if info != nil {
		return info, nil
	}

	return nil, fmt.Errorf("unable to get id %s", id)
}

func (m *memory) Insert(id string, info *common.UserInfo) error {
	return m.e.Insert(id, *info)
}

func (m *memory) Update(id string, info *common.UserInfo) error {
	return m.e.Insert(id, *info)
}

func (m *memory) Delete(id string) bool {
	return m.e.Delete(id)
}

func (m *memory) Exists(id string) bool {
	info := m.e.Get(id)
	if info != nil {
		return true
	}
	return false
}
