package memcache

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/HariniGB/login-provider/common"
	"github.com/HariniGB/login-provider/storage"
	"github.com/bradfitz/gomemcache/memcache"
)

type memcached struct {
	c       *memcache.Client
	timeout int32
}

func NewMemoryStorage(host, port string, timeout time.Duration) storage.Storage {
	client := memcache.New(fmt.Sprintf("%s:%s", host, port))

	return &memcached{c: client, timeout: int32(timeout.Seconds())}
}

func (m *memcached) Get(id string) (*common.UserInfo, error) {
	item, err := m.c.Get(id)

	if err != nil {
		log.Printf("Unable to get memcache key %s due to error %v\n", id, err)
		return nil, err
	}

	info := common.UserInfo{}
	err = json.Unmarshal(item.Value, &info)
	if err != nil {
		log.Printf("Unable to decode json for key %s due to error %v\n", id, err)
		return nil, err
	}

	return &info, nil
}

func (m *memcached) Insert(id string, info *common.UserInfo) error {
	bytes, err := json.Marshal(*info)
	if err != nil {
		return fmt.Errorf("Unable to marshal userinfo due to err: %v\n", err)
	}
	item := &memcache.Item{
		Key:        id,
		Value:      bytes,
		Expiration: m.timeout,
	}
	err = m.c.Set(item)
	if err != nil {
		log.Printf("Unable to insert memcache key %s due to error %v\n", id, err)
	}
	return err
}

func (m *memcached) Update(id string, info *common.UserInfo) error {
	bytes, err := json.Marshal(*info)
	if err != nil {
		return fmt.Errorf("Unable to marshal userinfo due to err: %v\n", err)
	}
	item := &memcache.Item{
		Key:        id,
		Value:      bytes,
		Expiration: m.timeout,
	}
	err = m.c.Replace(item)
	if err != nil {
		log.Printf("Unable to update memcache key %s due to error %v\n", id, err)
	}
	return err
}

func (m *memcached) Delete(id string) bool {
	err := m.c.Delete(id)
	if err != nil {
		log.Printf("Unable to delete memcache key %s due to error %v\n", id, err)
		return false
	}
	return true
}

func (m *memcached) Exists(id string) bool {
	info, _ := m.Get(id)
	if info != nil {
		return true
	}
	return false
}
