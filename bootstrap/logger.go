package bootstrap

import "desafio-itau-back-grpc/logger"

func InitLogger(debugMode bool, serviceName string) {
	logger.SetupLogging(debugMode, serviceName, logger.DEBUG)
	defer logger.AppLogger.Close()
}
