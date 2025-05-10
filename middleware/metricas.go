package middleware

import (
	"context"
	"desafio-itau-back-grpc/logger"
	"desafio-itau-back-grpc/modelos"
	"net"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

var (
	regMetricasGlobal        *modelos.Registro
	grpcServidorTratadoTotal *modelos.CounterVec
)

func SetupMetricas(registro *modelos.Registro, aplicacao string) {
	regMetricasGlobal = registro
	grpcServidorTratadoTotal = modelos.NovoCounterVec(
		registro,
		aplicacao,
		"Total de requisições gRPC tratadas pelo servidor.",
		[]string{"grpc_servico", "grpc_metodo", "grpc_codigo"},
	)
}

func MetricasUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	inicio := time.Now()

	var clientIP string = "desconhecido"
	p, ok := peer.FromContext(ctx)
	if ok {
		if tcpAddr, ok := p.Addr.(*net.TCPAddr); ok {
			clientIP = tcpAddr.IP.String()
		} else if unixAddr, ok := p.Addr.(*net.UnixAddr); ok {
			clientIP = unixAddr.Name + " (unix_socket)"
		} else {
			clientIP = p.Addr.String()
		}
	}

	partes := strings.Split(info.FullMethod, "/")
	servicoNome := "desconhecido"
	metodoNome := "desconhecido"
	if len(partes) == 3 {
		servicoNome = partes[1]
		metodoNome = partes[2]
	}
	resp, err := handler(ctx, req)

	duracao := time.Since(inicio)
	codigoStatus := status.Code(err).String()

	if grpcServidorTratadoTotal != nil {
		counterEspecifico := grpcServidorTratadoTotal.WithLabelValues(servicoNome, metodoNome, codigoStatus)
		counterEspecifico.Inc()
	} else {
		logger.AppLogger.Warn("grpcServidorTratadoTotal não inicializado!")
	}

	logger.AppLogger.Info("Log Interceptor: Client=%s  Servico=%s, Metodo=%s, Duracao=%s, Status=%s Em=%s\n",
		clientIP, servicoNome, metodoNome, duracao, codigoStatus, time.Now())

	return resp, err
}
