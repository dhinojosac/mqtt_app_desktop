package logger

import "go.uber.org/zap"

var logger *zap.SugaredLogger

//Init initialize unique instance of logger module
func Init() {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	logger = zapLogger.Sugar()
}

//Get is used to get the instance of logger module
func Get() *zap.SugaredLogger {
	return logger
}

//Close is used to close the instance of logger module
func Close() {
	logger.Sync()
}
