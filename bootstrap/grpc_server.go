package bootstrap

import (
	"desafio-itau-back-grpc/logger"
	"net"
)

func InitGRPCListener(addr string) (net.Listener, string, error) {
	if addr == "" {
		addr = ":4044"
		logger.AppLogger.Info("Variável váriavel de ambiente não definida não definida. Usando porta padrão: %s", addr)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, "", err
	}

	logger.AppLogger.Info("Servidor gRPC escutando em %s", addr)
	return listener, addr, nil
}
