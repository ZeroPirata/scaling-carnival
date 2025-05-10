package main

import (
	"desafio-itau-back-grpc/bootstrap"
	"desafio-itau-back-grpc/disk"
	"desafio-itau-back-grpc/logger"
	"desafio-itau-back-grpc/server"
	server_pb "desafio-itau-back-grpc/server/pb"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	debug := flag.Bool("debug", false, "Habilita logs no console (stdout/stderr) em vez de arquivo.")
	flag.Parse()

	bootstrap.InitLogger(*debug, "probes")
	bootstrap.LoadEnv()

	logger.AppLogger.Info("Iniciando serviços...")
	meuArmazem := disk.GetInstanciaUnica()
	defer logger.AppLogger.Close()
	client := server.Client(meuArmazem)
	server := server.Server(client)

	listener, addr, err := bootstrap.InitGRPCListener()
	if err != nil {
		logger.AppLogger.Fatal("Falha ao iniciar o listener: %v", err)
	}

	grpcServer := grpc.NewServer()
	server_pb.RegisterTransacaoServiceServer(grpcServer, server)

	errChan := make(chan error, 1)
	go func() {
		logger.AppLogger.Info("Iniciando gRPC em %s", addr)
		if err := grpcServer.Serve(listener); err != nil {
			errChan <- err
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-signalChan:
		logger.AppLogger.Warn("Sinal recebido: %v. Encerrando...", sig)
		grpcServer.GracefulStop()
		logger.AppLogger.Warn("Servidor finalizado com sucesso.")

	case err := <-errChan:
		logger.AppLogger.Error("Erro ao executar servidor: %v", err)
	}

}
