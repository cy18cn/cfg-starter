package cfg

import (
	"os"

	"github.com/spf13/viper"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func newZapLogger() (*zap.Logger, error) {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
	}

	level := zap.NewAtomicLevel()
	level.UnmarshalText([]byte(viper.GetString("logLevel")))
	// level := zap.NewAtomicLevelAt(zap.DebugLevel)

	config := zap.Config{
		Level:            level,                                                                     // 日志级别
		Development:      true,                                                                      // 开发模式，堆栈跟踪
		Encoding:         "json",                                                                    // 输出格式 console 或 json
		EncoderConfig:    encoderConfig,                                                             // 编码器配置
		InitialFields:    map[string]interface{}{"serviceName": viper.GetString("serviceName")}, // 初始化字段，如：添加一个服务器名称
		OutputPaths:      []string{"stdout", viper.GetString("logFile")},                        // 输出到指定文件 stdout（标准输出，正常颜色） stderr（错误输出，红色）
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}

func newLJZapLogger() *zap.Logger {
	hook := lumberjack.Logger{
		Filename:   viper.GetString("logFile"),
		MaxSize:    128,
		MaxBackups: 30,
		MaxAge:     30,
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
	}

	level := zap.NewAtomicLevel()
	level.UnmarshalText([]byte(viper.GetString("logLevel")))

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)),
		level,
	)

	caller := zap.AddCaller()
	dev := zap.Development()

	field := zap.Fields(zap.String("ServiceName", viper.GetString("serviceName")))
	logger := zap.New(core, caller, dev, field)

	return logger
}
