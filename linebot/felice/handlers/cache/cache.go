package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheHandler struct {
	cache *redis.Client
}

func NewCacheHandler() *CacheHandler {
	CACHE_URL := os.Getenv("CACHE_URL")
	client := redis.NewClient(
		&redis.Options{
			Addr:     CACHE_URL,
			Password: "",
			DB:       0,
		},
	)
	// test connection
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	return &CacheHandler{
		cache: client,
	}
}

func (h *CacheHandler) Get(ctx context.Context, key string, v interface{}) error {
	res, err := h.cache.Get(ctx, key).Result()
	if err == redis.Nil {
		return fmt.Errorf("key not found")
	} else if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(res), v); err != nil {
		return err
	}
	return nil
}

func (h *CacheHandler) Set(ctx context.Context, key string, v interface{}, ttl time.Duration) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return h.cache.Set(ctx, key, data, ttl).Err()
}

func (h *CacheHandler) Delete(ctx context.Context, key string) error {
	return h.cache.Del(ctx, key).Err()
}

// SetNX sets key to value if key does not exist. It returns true if the key was set, false if the key was not set.
func (h *CacheHandler) SetNX(ctx context.Context, key string, v interface{}, ttl time.Duration) (bool, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return false, err
	}
	return h.cache.SetNX(ctx, key, data, ttl).Result()
}

func (h *CacheHandler) Close() {
	h.cache.Close()
}

func (h *CacheHandler) HSet(ctx context.Context, key string, field string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return h.cache.HSet(ctx, key, field, data).Err()
}

func (h *CacheHandler) HGet(ctx context.Context, key string, field string, v interface{}) error {
	res, err := h.cache.HGet(ctx, key, field).Result()
	if err == redis.Nil {
		return fmt.Errorf("key not found")
	} else if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(res), v); err != nil {
		return err
	}
	return nil
}

func (h *CacheHandler) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return h.cache.Expire(ctx, key, ttl).Err()
}

func (h *CacheHandler) LPush(ctx context.Context, key string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return h.cache.LPush(ctx, key, data).Err()
}

func (h *CacheHandler) LRange(ctx context.Context, key string, start int64, stop int64, v interface{}) error {
	res, err := h.cache.LRange(ctx, key, start, stop).Result()
	if err == redis.Nil {
		return fmt.Errorf("key not found")
	} else if err != nil {
		return err
	}
	var values []interface{}
	for _, r := range res {
		var value interface{}
		if err := json.Unmarshal([]byte(r), &value); err != nil {
			return err
		}
		values = append(values, value)
	}
	reflect.ValueOf(v).Elem().Set(reflect.ValueOf(values))
	return nil
}

func (h *CacheHandler) LTrim(ctx context.Context, key string, start int64, stop int64) error {
	return h.cache.LTrim(ctx, key, start, stop).Err()
}
