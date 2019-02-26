package zlog

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var zLog *zap.Logger

//func init() {
//	zLog = newLJZapLogger()
//}

type LogOptions struct {
	Dev        bool
	Level      string
	AppName    string
	LogFile    string
	ErrLogFile string
	MaxSize    int
	MaxAge     int
	MaxBackups int
}

func NewProductionLogger(appName string) (*zap.Logger, error) {
	return NewLJLogger(&LogOptions{
		Dev:        false,
		Level:      "info",
		AppName:    appName,
		LogFile:    "/log/log.log",
		ErrLogFile: "/log/error.log",
		MaxSize:    128,
		MaxAge:     30,
		MaxBackups: 30,
	})
}

func NewLogger(opts *LogOptions) (*zap.Logger, error) {
	return newZapLogger(opts)
}

func NewLJLogger(opts *LogOptions) (*zap.Logger, error) {
	if opts.LogFile == "" || opts.ErrLogFile == "" {
		return nil, fmt.Errorf("logFile and errLogFile, one of them or both must be set")
	}
	return newLJZapLogger(opts), nil
}

func newZapLogger(opts *LogOptions) (*zap.Logger, error) {
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
	level.UnmarshalText([]byte(opts.Level))
	// level := zap.NewAtomicLevelAt(zap.DebugLevel)

	var outputPaths []string
	if opts.LogFile != "" {
		outputPaths = []string{"stdout", opts.LogFile}
	} else {
		outputPaths = []string{"stdout"}
	}

	var errOutputPaths []string
	if opts.ErrLogFile != "" {
		errOutputPaths = []string{"stderr", opts.ErrLogFile}
	} else {
		errOutputPaths = []string{"stderr"}
	}

	config := zap.Config{
		Level:            level,                                       // 日志级别
		Development:      true,                                        // 开发模式，堆栈跟踪
		Encoding:         "json",                                      // 输出格式 console 或 json
		EncoderConfig:    encoderConfig,                               // 编码器配置
		InitialFields:    map[string]interface{}{"APP": opts.AppName}, // 初始化字段，如：添加一个服务器名称
		OutputPaths:      outputPaths,                                 // 输出到指定文件 stdout（标准输出，正常颜色） stderr（错误输出，红色）
		ErrorOutputPaths: errOutputPaths,
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}

func newLJZapLogger(opts *LogOptions) *zap.Logger {
	hook := lumberjack.Logger{
		Filename:   opts.LogFile,
		MaxSize:    opts.MaxSize,
		MaxBackups: opts.MaxBackups,
		MaxAge:     opts.MaxAge,
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
	level.UnmarshalText([]byte(opts.Level))

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)),
		level,
	)

	caller := zap.AddCaller()
	dev := zap.Development()

	field := zap.Fields(zap.String("APP", opts.AppName))
	logger := zap.New(core, caller, dev, field)

	return logger
}
