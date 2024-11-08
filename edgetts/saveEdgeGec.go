package edgetts

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"go-edgetts/logger"
	redis2 "go-edgetts/redis"
	"time"
)

// redisSaveEdgeGec 函数用于将 GecInfo 对象保存到 Redis 中，并设置过期时间。
func redisSaveEdgeGec(Gec *GecInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	expireTime := Gec.Expiration
	now := time.Now().Unix()
	logger.SugarLogger.Debugf("SaveEdgeGec expireTime: %d, now: %d", expireTime, now)
	// 将 now 转换为 int64 后减去 expireTime
	ttl := time.Duration(expireTime - now)
	// 将 Gec 转换为 JSON
	gecData, err := json.Marshal(Gec)
	if err != nil {
		return err
	}
	key := "Gec"
	err = redis2.Rdb.Set(ctx, key, gecData, ttl*time.Second+180*time.Second).Err()
	return err
}

// 检查给定的GEC是否与从Redis获取的新GEC匹配
func checkGec(gec string) bool {
	newGec, err := redisGetEdgeGec()
	if err != nil {
		return false
	}
	if gec == newGec.SecMsgEC {
		return true
	}
	return false
}

func redisGetEdgeGec() (*GecInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	key := "Gec"

	// 从 Redis 中获取 JSON 字符串
	data, err := redis2.Rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		// 处理键为空的情况，可以返回一个自定义的错误或 nil 值
		return nil, errors.New("gec Key expired")
	} else if err != nil {
		return nil, err
	}

	// 将 JSON 字符串转换为结构体
	var Gec GecInfo
	if err := json.Unmarshal([]byte(data), &Gec); err != nil {
		return nil, err
	}
	return &Gec, nil
}
