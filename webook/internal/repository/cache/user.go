package cache

import (
	"basic-go/webook/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrKeyNotExist = redis.Nil

type UserCache interface {
	Get(ctx context.Context, userId string) (domain.UserProfile, error)
	Set(ctx context.Context, du domain.UserProfile) error
}

type RedisUserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func (c *RedisUserCache) Get(ctx context.Context, userId string) (domain.UserProfile, error) {
	key := c.key(userId)
	// 假定这个地方使用json序列化，然后可以使用给反序列化data
	data, err := c.cmd.Get(ctx, key).Result()
	if err != nil {
		return domain.UserProfile{}, err
	}
	var u domain.UserProfile
	err = json.Unmarshal([]byte(data), &u)
	return u, err
}

func (c *RedisUserCache) Set(ctx context.Context, du domain.UserProfile) error {
	key := c.key(du.User_id)
	data, err := json.Marshal(du)
	if err != nil {
		return err
	}
	return c.cmd.Set(ctx, key, data, c.expiration).Err()
}

func (c *RedisUserCache) key(userId string) string {
	return fmt.Sprintf("user:info:%d", userId)
}

func NewRedisUserCache(cmd redis.Cmdable) UserCache {
	return &RedisUserCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}
