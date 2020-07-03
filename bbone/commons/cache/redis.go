package cache

import (
	"time"

	"github.com/go-redis/cache/v7"
	"github.com/go-redis/redis/v7"
	"github.com/vmihailenco/msgpack/v4"
)

type rediscache struct {
	codec *cache.Codec
}

//RedisRing redis ring
type RedisRing map[string]string

//NewRedis new redis cache
func NewRedis(addrs RedisRing) (Cache, error) {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: addrs,
	})
	codec := &cache.Codec{
		Redis: ring,
		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}
	return &rediscache{
		codec: codec,
	}, nil
}

func (r *rediscache) SetWithExpTime(key string, value interface{}, exp time.Duration) error {
	return r.codec.Set(&cache.Item{Key: key, Object: value, Expiration: exp})
}

func (r *rediscache) Set(key string, value interface{}) error {
	exp := time.Minute * 10
	return r.SetWithExpTime(key, value, exp)
}

func (r *rediscache) Get(key string, obj interface{}) error {
	err := r.codec.Get(key, obj)
	return err
}

func (r *rediscache) Delete(key string) error {
	return r.codec.Delete(key)
}
