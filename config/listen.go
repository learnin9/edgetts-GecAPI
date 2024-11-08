package config // Package config 定义当前的API接口服务端口信息
import (
	"errors"
	"go-edgetts/logger"
	"gopkg.in/ini.v1"
	"log"
)

// GecEcAddress 定义 GecEcAddress 结构体
type GecEcAddress struct {
	Address string
}

// GecInfo 全局变量
var GecInfo *GecEcAddress

func ServerListen() string { // 服务器监听函数，读取配置文件中的端口
	cfg, err := ini.Load("server.conf")
	if err != nil {
		logger.SugarLogger.Fatal("配置文件读取错误", err)
	}
	ListenPort := cfg.Section("Server").Key("Port").String()
	return ":" + ListenPort
}
func GetGECAddress() (string, error) {
	cfg, err := ini.Load("server.conf")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 从配置中读取地址
	address := cfg.Section("edgeTTS").Key("updateGEC").String()
	if address == "" {
		return "", errors.New("GEC地址为空")
	}

	// 创建 GecEcAddress 实例并赋值
	GecInfo = &GecEcAddress{Address: address}
	return GecInfo.Address, nil
}
