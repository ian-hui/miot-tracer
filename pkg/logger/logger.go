package logger

import (
	mttypes "miot_tracing_go/mtTypes"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Miotlogger *zap.SugaredLogger

func init() {
	//清空日志文件的内容
	os.Truncate(mttypes.LogAddr, 0)
	//初始化日志
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	defaultLogLevel := zapcore.DebugLevel // 设置 loglevel，debug表示所有日志都输出，info表示只输出info以上的日志

	logFile, err := os.OpenFile(mttypes.LogAddr, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 06666)
	if err != nil {
		panic(err)
	}
	writer := zapcore.AddSync(logFile)

	logger := zap.New(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	Miotlogger = logger.Sugar()

}
