package main

import (
	"flag"
	"fmt"
	"github.com/gookit/color"
	"go-edgetts/Router"
	"go-edgetts/config"
	"go-edgetts/logger"
	"go-edgetts/redis"
	"log"
	"os"
)

var (
	version   string // 版本号
	buildTime string // 构建时间
)

func main() {
	// 如果传递了 --version 参数，则输出版本号和构建时间，并退出程序
	versionFlag := flag.Bool("version", false, "Print the version and build time")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Build Time: %s\n", buildTime)
		os.Exit(0)
	}

	ascii := `
:::::::::: :::::::::   ::::::::  :::::::::: ::::::::::: :::::::::::  ::::::::
:+:        :+:    :+: :+:    :+: :+:            :+:         :+:     :+:    :+:
+:+        +:+    +:+ +:+        +:+            +:+         +:+     +:+
+#++:++#   +#+    +:+ :#:        +#++:++#       +#+         +#+     +#++:++#++
+#+        +#+    +#+ +#+   +#+# +#+            +#+         +#+            +#+
#+#        #+#    #+# #+#    #+# #+#            #+#         #+#     #+#    #+#
########## #########   ########  ##########     ###         ###      ########
	`
	colors := []color.Color{
		color.FgRed,
		color.FgYellow,
		color.FgGreen,
		color.FgCyan,
		color.FgBlue,
		color.FgMagenta,
	}
	for i, c := range ascii {
		colors[i%len(colors)].Print(string(c))
	}
	println()

	logger.Init() // 初始化日志记录器
	_, err := config.GetGECAddress()
	if err != nil {
		log.Fatalf("Error getting GEC address: %v", err)
	}
	err = redis.ConnRedis()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}
	r := Router.SetupRouter()
	err = r.Run(config.ServerListen())
	if err != nil {
		logger.SugarLogger.Fatalf("Run err: %v", err)
	}
}
