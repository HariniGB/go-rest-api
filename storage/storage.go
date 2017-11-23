package storage

import "github.com/HariniGB/login-provider/common"

type Storage interface {
	Get(string) (*common.UserInfo, error)
	Insert(string, *common.UserInfo) error
	Update(string, *common.UserInfo) error
	Delete(string) bool
	Exists(string) bool
}
