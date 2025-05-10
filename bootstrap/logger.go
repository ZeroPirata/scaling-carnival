package bootstrap

import "desafio-itau-back-grpc/logger"

func InitLogger(debugMode bool, serviceName string) {
	var minLevel logger.Level
	if debugMode {
		minLevel = logger.DEBUG
	} else {
		minLevel = logger.INFO
	}
	logger.SetupLogging(debugMode, serviceName, minLevel)
}
