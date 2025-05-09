package server

import (
	"context"
	"desafio-itau-back-grpc/logger"
	server_pb "desafio-itau-back-grpc/server/pb"

	"google.golang.org/protobuf/types/known/emptypb"
)

type GrpcServer struct {
	server_pb.UnimplementedTransacaoServiceServer
	client *ClientService
}

func Server(client *ClientService) *GrpcServer {
	if client == nil {
		logger.AppLogger.Error("Tentativa de criar grpc com o cliente foi nula")
	}
	return &GrpcServer{
		client: client,
	}
}

func (s *GrpcServer) CriarTransacao(ctx context.Context, params *server_pb.CriarTransacaoRequest) (res *emptypb.Empty, err error) {
	res, err = s.client.CriarTransacao(ctx, params)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (s *GrpcServer) GetEstatistica(ctx context.Context, params *emptypb.Empty) (res *server_pb.EstatisticaResponse, err error) {
	res, err = s.client.GetEstatistica(ctx, params)
	if err != nil {
		return nil, err
	}
	return res, err
}
