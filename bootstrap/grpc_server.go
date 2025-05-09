package bootstrap

import (
	"desafio-itau-back-grpc/logger"
	"net"
	"os"
)

func InitGRPCListener() (net.Listener, string, error) {
	addr := os.Getenv("GRPC_ADDR_POSTGRES")
	if addr == "" {
		addr = ":4044"
		logger.AppLogger.Info("Variável GRPC_ADDR_POSTGRES não definida. Usando porta padrão: %s", addr)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, "", err
	}

	logger.AppLogger.Info("Servidor gRPC escutando em %s", addr)
	return listener, addr, nil
}
