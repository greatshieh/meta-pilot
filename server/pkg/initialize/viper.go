package initialize

import (
	"flag"
	"fmt"
	"os"

	"server/pkg/global"
	"server/pkg/initialize/internal"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Viper 初始化并返回viper实例，用于加载和管理配置文件
//
// 参数:
//
//	inputPath: 可选参数，指定配置文件路径。如果提供，将优先使用该路径
//
// 返回值:
//
//	*viper.Viper: 配置管理器实例，包含已加载的配置信息
//
// 配置文件加载优先级（从高到低）:
// 1. 函数参数 inputPath[0]
// 2. 命令行参数 -c
// 3. 环境变量 CONFIG
// 4. 根据gin模式自动选择默认配置文件
//   - Debug模式: config.develop.yaml
//   - Release模式: config.yaml
//   - Test模式: config.test.yaml
//
// 功能特点:
// - 支持配置文件热重载
// - 自动将配置解析到全局变量 global.MPA_CONFIG
// - 提供详细的日志输出，显示当前使用的配置来源和路径
func Viper(inputPath ...string) *viper.Viper {
	var config string // 配置文件路径变量

	// 检查是否通过函数参数指定了配置文件路径
	if len(inputPath) == 0 {
		// 定义命令行参数 -c，用于指定配置文件
		flag.StringVar(&config, "c", "", "choose config file.")
		flag.Parse() // 解析命令行参数

		// 检查命令行参数是否为空
		if config == "" { // 判断命令行参数是否为空, 优先级: 命令行 > 环境变量 > 默认值
			// 检查是否设置了环境变量 CONFIG
			if configEnv := os.Getenv(internal.ConfigEnv); configEnv == "" { // 判断 internal.ConfigEnv 常量存储的环境变量是否为空
				// 根据gin运行模式选择默认配置文件
				switch gin.Mode() {
				case gin.DebugMode:
					config = internal.ConfigDefaultFile // 使用开发环境配置文件
					fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.EnvGinMode, internal.ConfigDefaultFile)
				case gin.ReleaseMode:
					config = internal.ConfigReleaseFile // 使用生产环境配置文件
					fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.EnvGinMode, internal.ConfigReleaseFile)
				case gin.TestMode:
					config = internal.ConfigTestFile // 使用测试环境配置文件
					fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.EnvGinMode, internal.ConfigTestFile)
				}
			} else {
				// 使用环境变量指定的配置文件路径
				config = configEnv
				fmt.Printf("您正在使用%s环境变量,config的路径为%s\n", internal.ConfigEnv, config)
			}
		} else {
			// 使用命令行参数指定的配置文件路径
			fmt.Printf("您正在使用命令行的-c参数传递的值,config的路径为%s\n", config)
		}
	} else { // 函数传递的可变参数的第一个值赋值于config
		// 使用函数参数指定的配置文件路径
		config = inputPath[0]
		fmt.Printf("您正在使用func Viper()传递的值,config的路径为%s\n", config)
	}

	// 创建新的viper实例
	v := viper.New()
	v.SetConfigFile(config) // 设置配置文件路径
	v.SetConfigType("yaml") // 设置配置文件类型为yaml

	// 读取配置文件
	err := v.ReadInConfig()
	if err != nil {
		// 配置文件读取失败时，抛出异常终止程序
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	// 启用配置文件监听功能，当配置文件变化时自动重新加载
	v.WatchConfig()

	// 设置配置文件变化回调函数
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name) // 输出配置文件变化信息
		// 将新配置解析到全局变量
		if err = v.Unmarshal(&global.MPA_CONFIG); err != nil {
			fmt.Println(err) // 解析失败时输出错误信息
		}
	})

	// 将初始配置解析到全局变量
	if err = v.Unmarshal(&global.MPA_CONFIG); err != nil {
		panic(err) // 解析失败时抛出异常终止程序
	}

	return v // 返回配置管理器实例
}
