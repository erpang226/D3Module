package config

import (
	"flag"
	"k8s.io/klog/v2"
)

type LogConfig struct {
	Level     string `yaml:"level"` // 日志级别
	LogRotate struct {
		Filename   string `yaml:"fileName"`   // 日志文件名
		MaxSize    string `yaml:"maxSize"`    // 日志文件最大大小（MB）
		MaxBackups int    `yaml:"maxBackups"` // 保留的日志文件数量
		MaxAge     int    `yaml:"maxAge"`     // 日志文件最长存活时间（天）
		Compress   bool   `yaml:"compress"`   // 是否压缩日志文件
	} `yaml:"logRotate"` // 日志轮转配置
}

// LogInit 初始化日志
func (logConfig *LogConfig) LogInit() error {
	klog.InitFlags(nil)
	if len(logConfig.LogRotate.Filename) > 0 {
		err := flag.Set("logtostderr", "false")
		if err != nil {
			return err
		}
		err = flag.Set("log_file", logConfig.LogRotate.Filename)
		if err != nil {
			return err
		}
		err = flag.Set("log_file_max_size", logConfig.LogRotate.MaxSize)
		if err != nil {
			return err
		}
	}
	err := flag.Set("v", logConfig.Level)
	if err != nil {
		return err
	}

	flag.Parse()

	return nil
}

func (logConfig *LogConfig) FlushLogs() {
	klog.Flush()
}
