package zlog

import "go.uber.org/zap"

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Info(msg string, fields ...zap.Field) {
	zLog.Info(msg, fields...)
}

// Debug logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Debug(msg string, fields ...zap.Field) {
	zLog.Debug(msg, fields...)
}

// Warn logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Warn(msg string, fields ...zap.Field) {
	zLog.Warn(msg, fields...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func Fatal(msg string, fields ...zap.Field) {
	zLog.Fatal(msg, fields...)
}

// Error logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Error(msg string, fields ...zap.Field) {
	zLog.Error(msg, fields...)
}

// Sync calls the underlying Core's Sync method, flushing any buffered log
// entries. Applications should take care to call Sync before exiting.
func Sync() error {
	return zLog.Sync()
}

// Infof uses fmt.Sprintf to log a templated message.
func Infof(template string, args ...interface{}) {
	zLog.Sugar().Infof(template, args)
}

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(template string, args ...interface{}) {
	zLog.Sugar().Debugf(template, args)
}

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(template string, args ...interface{}) {
	zLog.Sugar().Warnf(template, args)
}

// Fatalf uses fmt.Sprintf to log a templated message.
func Fatalf(template string, args ...interface{}) {
	zLog.Sugar().Fatalf(template, args)
}

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(template string, args ...interface{}) {
	zLog.Sugar().Errorf(template, args)
}

func GetLogger() *zap.Logger {
	return zLog
}

func InitProduction(appName string) (err error) {
	zLog, err = NewProductionLogger(appName)
	return err
}
