package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/coocood/freecache"
	"github.com/redis/go-redis/v9"
	"strconv"
	"sync"
)

var (
	//go:embed lua/set_code.lua
	luaSetCode string
	//go:embed lua/verify_code.lua
	luaVerifyCode string

	ErrCodeSendTooMany   = errors.New("发送太频繁")
	ErrCodeVerifyTooMany = errors.New("验证次数太多")
	ErrNotFount          = freecache.ErrNotFound
)

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inCode string) (bool, error)
}

// RedisCodeCache ----------------------------------------------使用redis缓存实现--------------------------------------------//
type RedisCodeCache struct {
	cmd redis.Cmdable
}

func NewRedisCodeCache(cmd redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		cmd: cmd,
	}
}

func (c *RedisCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := c.cmd.Eval(ctx, luaSetCode, []string{c.Key(biz, phone)}, code).Int()
	if err != nil {
		// 调用redis出了问题
		return err
	}
	switch res {
	case -2:
		return errors.New("验证码存在，但是没有过期时间")
	case -1:
		return ErrCodeSendTooMany
	default:
		return nil
	}
}

func (c *RedisCodeCache) Verify(ctx context.Context, biz, phone, inCode string) (bool, error) {
	res, err := c.cmd.Eval(ctx, luaVerifyCode, []string{c.Key(biz, phone)}, inCode).Int()
	if err != nil {
		// 调用redis出了问题
		return false, err
	}
	switch res {
	case -2:
		// 校验验证码错误
		return false, nil
	case -1:
		// 验证次数耗尽，意味着验证频繁
		return false, ErrCodeVerifyTooMany
	default:
		// 验证成功
		return true, nil
	}
}

func (c *RedisCodeCache) Key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

// MemoryCodeCache ---------------------------------------使用本地缓存实现(freeCache)--------------------------------------//
type MemoryCodeCache struct {
	cc *freecache.Cache
}

func NewMemoryCodeCache(cache *freecache.Cache) CodeCache {
	return &MemoryCodeCache{
		cc: cache,
	}
}

func (c *MemoryCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := c.setWithCache(c.Key(biz, phone), code)
	if err != nil {
		// 调用redis出了问题
		return err
	}
	switch res {
	case -2:
		return errors.New("验证码存在，但是没有过期时间")
	case -1:
		return ErrCodeSendTooMany
	default:
		return nil
	}
}

func (c *MemoryCodeCache) Verify(ctx context.Context, biz, phone, inCode string) (bool, error) {
	res, err := c.VerifyWithCache(c.Key(biz, phone), inCode)
	if err != nil {
		// 调用redis出了问题
		return false, err
	}
	switch res {
	case -2:
		// 校验验证码错误
		return false, nil
	case -1:
		// 验证次数耗尽，意味着验证频繁
		return false, ErrCodeVerifyTooMany
	default:
		// 验证成功
		return true, nil
	}
}

func (c *MemoryCodeCache) Key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

func (c *MemoryCodeCache) setWithCache(key, code string) (int32, error) {
	var lock sync.Mutex
	lock.Lock()
	defer lock.Unlock()
	// 拼接字符串，cntKey表示可以重试的次数
	cntKey := key + ":cnt"
	// 获取验证码的有效时间
	lastTime, err := c.cc.TTL([]byte(key))
	if err != nil && !errors.Is(err, ErrNotFount) {
		return -2, err
	} else if lastTime < 540 || errors.Is(err, ErrNotFount) {
		// 插入验证码和重试次数
		c.cc.Set([]byte(key), []byte(code), 600)
		c.cc.Set([]byte(cntKey), []byte(strconv.Itoa(3)), 600)
		return 0, nil
	} else {
		// 获取验证码太频繁
		return -1, nil
	}
}

func (c *MemoryCodeCache) VerifyWithCache(key, incode string) (int32, error) {
	var lock sync.Mutex
	lock.Lock()
	defer lock.Unlock()

	// 获得缓存中的验证码
	code, err := c.cc.Get([]byte(key))
	if err != nil {
		return -1, err
	}
	// 拼接字符串，cntKey表示可以重试的次数
	cntKey := key + ":cnt"
	cntStr, err := c.cc.Get([]byte(cntKey))
	if err != nil {
		return -1, err
	}
	// 字符串转数字
	cnt, err := strconv.Atoi(string(cntStr))
	if err != nil {
		return -1, err
	}
	// 判断是否剩余重试次数
	if cnt <= 0 {
		// 验证次数耗尽了
		return -1, nil
	}
	if string(code) == incode {
		// 相等，次数置零
		c.cc.Set([]byte(cntKey), []byte(strconv.Itoa(0)), -1)
		return 0, nil
	} else {
		// 不相等，用户输入验证码错误，次数减一
		c.cc.Set([]byte(cntKey), []byte(strconv.Itoa(cnt-1)), 600)
		return -2, nil
	}
}
