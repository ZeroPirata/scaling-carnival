package main

import (
	"desafio-itau-back-grpc/bootstrap"
	"desafio-itau-back-grpc/disk"
	"desafio-itau-back-grpc/logger"
	"desafio-itau-back-grpc/middleware"
	"desafio-itau-back-grpc/modelos"
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

	bootstrap.InitLogger(*debug, "itau-desafio-grpc")
	bootstrap.LoadEnv()

	logger.AppLogger.Info("Iniciando serviços...")
	registroDeMetricas := modelos.NovoRegistro()

	middleware.SetupMetricas(registroDeMetricas, "itau-desafio-grpc")

	portaMetricas := os.Getenv("METRICS_PORT")
	if portaMetricas == "" {
		portaMetricas = "9091"
	}
	go server.IniciarServidorMetricas(":"+portaMetricas, registroDeMetricas)

	meuArmazem := disk.GetInstanciaUnica()
	defer logger.AppLogger.Close()

	client := server.Client(meuArmazem)
	server := server.Server(client)

	port := os.Getenv("PORT")
	listener, addr, err := bootstrap.InitGRPCListener(port)
	if err != nil {
		logger.AppLogger.Fatal("Falha ao iniciar o listener: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.MetricasUnaryInterceptor),
	)
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
