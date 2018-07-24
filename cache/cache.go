package cache

import "time"

//Cache interface
type Cache interface {
	Get(key string) interface{}
	GetString(key string) string
	Set(key string, val interface{}, timeout time.Duration) error
	SetString(key string, val string, timeout time.Duration) error
	IsExist(key string) bool
	Delete(key string) error
}
