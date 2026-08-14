package initialize

import (
	"fmt"
	"os"

	"server/pkg/global"
	"server/pkg/initialize/internal"
	"server/pkg/utils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Zap 获取 zap.Logger

func Zap() (logger *zap.Logger) {
	if ok, _ := utils.PathExists(global.MPA_CONFIG.Zap.Director); !ok { // 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", global.MPA_CONFIG.Zap.Director)
		_ = os.Mkdir(global.MPA_CONFIG.Zap.Director, os.ModePerm)
	}

	cores := internal.Zap.GetZapCores()
	logger = zap.New(zapcore.NewTee(cores...))

	if global.MPA_CONFIG.Zap.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger
}
