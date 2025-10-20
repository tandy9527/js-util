package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/tandy9527/js-util/logger"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	rdb     *redis.Client
	timeout time.Duration
}

var (
	redisMap = make(map[string]*RedisClient)
	once     sync.Once
)

// LoadRedis 初始化多个 Redis 客户端
func LoadRedis(path string) {
	cfg := LoadRedisConf(path)
	once.Do(func() {
		// 配置的多个DB
		for name, c := range cfg.Redis {
			rdb := redis.NewClient(&redis.Options{
				Addr:         c.Addr,
				Password:     c.Password,
				DB:           c.DB,
				PoolSize:     c.PoolSize,
				MinIdleConns: c.MinIdleConns,
				PoolTimeout:  time.Duration(c.PoolTimeout) * time.Second,
			})

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			if err := rdb.Ping(ctx).Err(); err != nil {
				panic(fmt.Sprintf("Redis[%s] Connection failed: %v", name, err))
			}

			redisMap[name] = &RedisClient{
				rdb:     rdb,
				timeout: 2 * time.Second,
			}
			logger.Infof("Redis[%s] Connection successful", name)
		}
	})
}

// CloseRedis 关闭所有 Redis
func CloseRedis() {
	for name, cli := range redisMap {
		if cli != nil {
			if err := cli.rdb.Close(); err != nil {
				logger.Errorf("Redis[%s] Close failed: %v", name, err)
			} else {
				logger.Infof("Redis[%s] Closed", name)
			}
		}
	}
}

// GetDB 获取指定 Redis 客户端
func GetDB(name string) *RedisClient {
	return redisMap[name]
}

// 内部统一执行函数
func (c *RedisClient) do(fn func(ctx context.Context) (any, error), timeout ...time.Duration) (any, error) {
	t := c.timeout
	if len(timeout) > 0 {
		t = timeout[0]
	}
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()
	//start := time.Now()
	res, err := fn(ctx)
	//log.Printf("[Redis] cost=%v err=%v", time.Since(start), err)
	return res, err
}

func (c *RedisClient) Expire(key string, expiration time.Duration, timeout ...time.Duration) error {
	_, err := c.do(func(ctx context.Context) (any, error) {
		return nil, c.rdb.Expire(ctx, key, expiration).Err() // (bool, error)
	}, timeout...)
	return err
}

// ---------------- 基础操作 ----------------

func (c *RedisClient) Set(key string, value any, expiration time.Duration, timeout ...time.Duration) error {
	_, err := c.do(func(ctx context.Context) (any, error) {
		return nil, c.rdb.Set(ctx, key, value, expiration).Err()
	}, timeout...)
	return err
}

func (c *RedisClient) Get(key string, timeout ...time.Duration) (string, error) {
	res, err := c.do(func(ctx context.Context) (any, error) {
		return c.rdb.Get(ctx, key).Result()
	}, timeout...)
	if err != nil {
		return "", err
	}
	return res.(string), nil
}

func (c *RedisClient) Exists(key string, timeout ...time.Duration) (bool, error) {
	res, err := c.do(func(ctx context.Context) (any, error) {
		return c.rdb.Exists(ctx, key).Result()
	}, timeout...)
	if err != nil {
		return false, err
	}
	if val, ok := res.(int64); ok {
		return val > 0, nil
	}
	return false, nil
}

func (c *RedisClient) Del(key string, timeout ...time.Duration) error {
	_, err := c.do(func(ctx context.Context) (any, error) {
		return nil, c.rdb.Del(ctx, key).Err()
	}, timeout...)
	return err
}

func (c *RedisClient) Incr(key string, timeout ...time.Duration) (int64, error) {
	res, err := c.do(func(ctx context.Context) (any, error) {
		return c.rdb.Incr(ctx, key).Result()
	}, timeout...)
	return res.(int64), err
}

func (c *RedisClient) Decr(key string, timeout ...time.Duration) (int64, error) {
	res, err := c.do(func(ctx context.Context) (any, error) {
		return c.rdb.Decr(ctx, key).Result()
	}, timeout...)
	return res.(int64), err
}

func (c *RedisClient) DecrBy(key string, decrement int64, timeout ...time.Duration) (int64, error) {
	res, err := c.do(func(ctx context.Context) (any, error) {
		return c.rdb.DecrBy(ctx, key, decrement).Result()
	}, timeout...)
	return res.(int64), err
}

// ---------------- Hash 操作 ----------------

// HSet 设置 hash key 字段和值
func (c *RedisClient) HSet(hashKey string, field string, value any, timeout ...time.Duration) error {
	_, err := c.do(func(ctx context.Context) (any, error) {
		return nil, c.rdb.HSet(ctx, hashKey, field, value).Err()
	}, timeout...)
	return err
}

// HMSet 批量设置 hash key 的多个字段和值
// 注意：v9 版本已不再使用 HMSET，直接用 HSet 替代
func (c *RedisClient) HMSet(key string, data map[string]any, timeout ...time.Duration) error {
	_, err := c.do(func(ctx context.Context) (any, error) {
		return c.rdb.HSet(ctx, key, data).Result()
	}, timeout...)
	return err
}

// HGet 获取 hash key 指定字段的值
func (c *RedisClient) HGet(hashKey string, field string, timeout ...time.Duration) (string, error) {
	res, err := c.do(func(ctx context.Context) (any, error) {
		return c.rdb.HGet(ctx, hashKey, field).Result()
	}, timeout...)
	if err != nil {
		return "", err
	}
	return res.(string), nil
}

// HDel 删除 hash key 指定字段
func (c *RedisClient) HDel(hashKey string, fields ...string) error {
	_, err := c.do(func(ctx context.Context) (any, error) {
		return nil, c.rdb.HDel(ctx, hashKey, fields...).Err()
	})
	return err
}

// HGetAll 获取整个 hash
func (c *RedisClient) HGetAll(hashKey string) (map[string]string, error) {
	res, err := c.do(func(ctx context.Context) (any, error) {
		return c.rdb.HGetAll(ctx, hashKey).Result()
	})
	if err != nil {
		return nil, err
	}
	return res.(map[string]string), nil
}

// --------------------------- Set 操作 ---------------------------

// SAdd 添加一个或多个成员到集合
func (c *RedisClient) SAdd(key string, members ...any) (int64, error) {
	res, err := c.do(func(ctx context.Context) (any, error) {
		return c.rdb.SAdd(ctx, key, members...).Result()
	})
	if err != nil {
		return 0, err
	}
	return res.(int64), nil
}

// SMembers 获取集合所有成员
func (c *RedisClient) SMembers(key string) ([]string, error) {
	res, err := c.do(func(ctx context.Context) (any, error) {
		return c.rdb.SMembers(ctx, key).Result()
	})
	if err != nil {
		return nil, err
	}
	return res.([]string), nil
}

// SIsMember 判断元素是否存在
func (c *RedisClient) SIsMember(key string, member any) (bool, error) {
	res, err := c.do(func(ctx context.Context) (any, error) {
		return c.rdb.SIsMember(ctx, key, member).Result()
	})
	if err != nil {
		return false, err
	}
	return res.(bool), nil
}

// SRem 删除集合中的一个或多个成员
func (c *RedisClient) SRem(key string, members ...any) (int64, error) {
	res, err := c.do(func(ctx context.Context) (any, error) {
		return c.rdb.SRem(ctx, key, members...).Result()
	})
	if err != nil {
		return 0, err
	}
	return res.(int64), nil
}

// SCard 获取集合大小
func (c *RedisClient) SCard(key string) (int64, error) {
	res, err := c.do(func(ctx context.Context) (any, error) {
		return c.rdb.SCard(ctx, key).Result()
	})
	if err != nil {
		return 0, err
	}
	return res.(int64), nil
}

// SPop 取出一个元素并移除
func (c *RedisClient) SPop(key string) (string, error) {
	res, err := c.do(func(ctx context.Context) (any, error) {
		return c.rdb.SPop(ctx, key).Result()
	})
	if err != nil {
		if err == redis.Nil { // 集合为空
			return "", nil
		}
		return "", err
	}
	return res.(string), nil
}

// ---------------- Lua 原子操作 ----------------
func (c *RedisClient) ExecLua(script string, keys []string, args ...any) (any, error) {
	return c.do(func(ctx context.Context) (any, error) {
		lua := redis.NewScript(script)
		return lua.Run(ctx, c.rdb, keys, args...).Result()
	})
}

// BRPopLPush 封装
func (c *RedisClient) BRPopLPush(source, dest string, timeoutSeconds int, timeout ...time.Duration) (string, error) {
	res, err := c.do(func(ctx context.Context) (any, error) {
		// 这里使用 Redis 原生 BRPopLPush
		return c.rdb.BRPopLPush(ctx, source, dest, time.Duration(timeoutSeconds)*time.Second).Result()
	}, timeout...)
	if err != nil {
		return "", err
	}

	if str, ok := res.(string); ok {
		return str, nil
	}
	return "", nil
}

// LRem 从 list 中删除指定的元素
func (c *RedisClient) LRem(key string, count int64, value any, timeout ...time.Duration) (int64, error) {
	result, err := c.do(func(ctx context.Context) (any, error) {
		return c.rdb.LRem(ctx, key, count, value).Result()
	}, timeout...)
	if err != nil {
		return 0, err
	}
	return result.(int64), nil
}

// LPush 将一个或多个值推入 list 左侧
func (c *RedisClient) LPush(key string, values ...any) (int64, error) {
	result, err := c.do(func(ctx context.Context) (any, error) {
		return c.rdb.LPush(ctx, key, values...).Result()
	})
	if err != nil {
		return 0, err
	}
	return result.(int64), nil
}

// RPop 从 Redis 队列右侧弹出一个元素
func (c *RedisClient) RPop(key string, timeout ...time.Duration) (string, error) {
	res, err := c.do(func(ctx context.Context) (any, error) {
		return c.rdb.RPop(ctx, key).Result()
	}, timeout...)
	if err != nil {
		return "", err
	}

	// 类型断言为 string
	if str, ok := res.(string); ok {
		return str, nil
	}
	return "", nil
}

// ----------------------------zset--------------------------------
// 添加或更新成员分数
func (c *RedisClient) ZAdd(key string, members ...redis.Z) (int64, error) {
	res, err := c.do(func(ctx context.Context) (any, error) {
		return c.rdb.ZAdd(ctx, key, members...).Result()
	})
	if err != nil {
		return 0, err
	}
	return res.(int64), nil
}

// 增加分数（支持正负值）
func (c *RedisClient) ZIncrBy(key string, increment float64, member string) (float64, error) {
	res, err := c.do(func(ctx context.Context) (any, error) {
		return c.rdb.ZIncrBy(ctx, key, increment, member).Result()
	})
	if err != nil {
		return 0, err
	}
	return res.(float64), nil
}

// 获取成员分数
func (c *RedisClient) ZScore(key, member string) (float64, error) {
	res, err := c.do(func(ctx context.Context) (any, error) {
		return c.rdb.ZScore(ctx, key, member).Result()
	})
	if err != nil {
		return 0, err
	}
	return res.(float64), nil
}

// 按分数从高到低取范围
func (c *RedisClient) ZRevRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	res, err := c.do(func(ctx context.Context) (any, error) {
		return c.rdb.ZRevRangeWithScores(ctx, key, start, stop).Result()
	})
	if err != nil {
		return nil, err
	}
	return res.([]redis.Z), nil
}

// 删除成员
func (c *RedisClient) ZRem(key string, members ...any) (int64, error) {
	res, err := c.do(func(ctx context.Context) (any, error) {
		return c.rdb.ZRem(ctx, key, members...).Result()
	})
	if err != nil {
		return 0, err
	}
	return res.(int64), nil
}

// 获取总成员数量
func (c *RedisClient) ZCard(key string) (int64, error) {
	res, err := c.do(func(ctx context.Context) (any, error) {
		return c.rdb.ZCard(ctx, key).Result()
	})
	if err != nil {
		return 0, err
	}
	return res.(int64), nil
}

// 删除分数区间内的成员
func (c *RedisClient) ZRemRangeByScore(key, min, max string) (int64, error) {
	res, err := c.do(func(ctx context.Context) (any, error) {
		return c.rdb.ZRemRangeByScore(ctx, key, min, max).Result()
	})
	if err != nil {
		return 0, err
	}
	return res.(int64), nil
}
