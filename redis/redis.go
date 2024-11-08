package redis

import (
	"context"
	"github.com/gookit/color"
	"github.com/redis/go-redis/v9"
	"go-edgetts/logger"
	"gopkg.in/ini.v1"
	"strings"
	"time"
)

var Rdb redis.UniversalClient

type Info struct {
	Addr string //地址信息
	Port string //端口
	Pass string //密码
	Mode string //模式是single还是cluster
}

func InfoRedis() Info {
	cfg, err := ini.Load("server.conf")
	if err != nil {
		logger.SugarLogger.Fatalf("读取文件错误:%s", err)
	}
	redisAddr := cfg.Section("Redis").Key("Address").String()
	redisPort := cfg.Section("Redis").Key("Port").String()
	auth := cfg.Section("Redis").Key("Auth").String()
	mode := cfg.Section("Redis").Key("Mode").String()
	redisTable := Info{
		Addr: redisAddr,
		Port: redisPort,
		Pass: auth,
		Mode: mode,
	}
	return redisTable
}

// ConnRedis Redis数据库初始化连接
func ConnRedis() (err error) {
	redisConfig := InfoRedis()
	if redisConfig.Mode == "cluster" {
		// Parse the comma-separated list of addresses for cluster mode
		addrs := strings.Split(redisConfig.Addr, ",")
		for i, addr := range addrs {
			addrs[i] = addr + ":" + redisConfig.Port
			logger.SugarLogger.Debugf("集群获取地址:%s", addrs[i])
		}
		Rdb = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:           addrs,
			Password:        redisConfig.Pass,
			PoolSize:        100, // Connection pool size
			MinIdleConns:    10,
			DialTimeout:     5 * time.Second,
			ReadTimeout:     3 * time.Second,
			WriteTimeout:    3 * time.Second,
			PoolTimeout:     4 * time.Second,
			MaxRetries:      0,
			MinRetryBackoff: 8 * time.Millisecond,
			MaxRetryBackoff: 512 * time.Millisecond,
		})
	} else { // 如果是single模式就是走单机模式
		Rdb = redis.NewClient(&redis.Options{
			Addr:            redisConfig.Addr + ":" + redisConfig.Port,
			Password:        redisConfig.Pass,
			DB:              0,   // Use default DB
			PoolSize:        100, // Connection pool size
			MinIdleConns:    10,
			DialTimeout:     5 * time.Second,
			ReadTimeout:     3 * time.Second,
			WriteTimeout:    3 * time.Second,
			PoolTimeout:     4 * time.Second,
			MaxRetries:      0,
			MinRetryBackoff: 8 * time.Millisecond,
			MaxRetryBackoff: 512 * time.Millisecond,
		})
	}

	// Using context to handle timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err = Rdb.Ping(ctx).Result(); err != nil {
		logger.SugarLogger.Warnf("Redis连接失败,错误原因:%s", err)
		color.Red.Println("[x]连接Redis数据库失败")
		return err
	} else {
		color.Green.Println("[✔]连接Redis数据库成功")
		logger.SugarLogger.Info("[✔]连接Redis数据库成功")
		return nil
	}
}
