package cache

import (
	"time"
)

//Cache ...
type Cache interface {
	Set(key string, value interface{}) error
	SetWithExpTime(key string, value interface{}, exp time.Duration) error
	Get(key string, obj interface{}) error
	Delete(key string) error
}
