package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var Logger *zap.Logger

// InitZap 初始化zap logger
func InitZap() *zap.Logger {
	Logger = zap.New(getEncoderCore(), zap.AddCaller())
	Logger.Info("logger init success")
	return Logger
}

func getEncoderCore() zapcore.Core {
	core := zapcore.NewCore(getEncoder(), getWriteSyncer(), zap.DebugLevel)
	return core
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}

func getEncoderConfig() zapcore.EncoderConfig {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = CustomTimeEncoder
	return config
}

func getWriteSyncer() zapcore.WriteSyncer {
	//file, err := os.OpenFile("./info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	//if err != nil {
	//	return nil
	//}
	return zapcore.AddSync(os.Stdout)
}

// CustomTimeEncoder 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006/01/02 15:04:05.000"))
}
