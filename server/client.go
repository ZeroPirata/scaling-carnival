package server

import (
	"context"
	"desafio-itau-back-grpc/disk"
	"desafio-itau-back-grpc/logger"
	"desafio-itau-back-grpc/modelos"
	server_pb "desafio-itau-back-grpc/server/pb"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ClientService struct {
	*disk.Armazem
}

func Client(disk *disk.Armazem) *ClientService {
	return &ClientService{
		Armazem: disk,
	}
}

func (s *ClientService) CriarTransacao(ctx context.Context, params *server_pb.CriarTransacaoRequest) (*emptypb.Empty, error) {
	transacaoPayload := params.GetTransacao()
	if transacaoPayload == nil {
		logger.AppLogger.Warn("Requisição de transação recebida sem payload de transação.")
		return nil, status.Errorf(codes.InvalidArgument, "payload da transação não fornecido")
	}
	valorTransacao := transacaoPayload.GetValor()
	dataHoraString := transacaoPayload.GetDataHora()

	if valorTransacao <= 0 {
		msgErro := fmt.Sprintf("valor da transação deve ser positivo, recebido: %.2f", valorTransacao)
		logger.AppLogger.Warn(msgErro)
		return nil, status.Errorf(codes.InvalidArgument, "valor da transação deve ser positivo")
	}

	if strings.TrimSpace(dataHoraString) == "" {
		logger.AppLogger.Warn("dataHora da transação (string) não foi fornecida ou está vazia.")
		return nil, status.Errorf(codes.InvalidArgument, "dataHora da transação deve ser fornecida")
	}

	dataHoraTransacao, err := time.Parse(time.RFC3339Nano, dataHoraString)
	if err != nil {
		msgErro := fmt.Sprintf("dataHora da transação ('%s') não está no formato ISO 8601 (RFC3339Nano) esperado. Erro de parse: %v", dataHoraString, err)
		logger.AppLogger.Warn(msgErro)
		return nil, status.Errorf(codes.InvalidArgument, "dataHora da transação com formato inválido. Esperado formato ISO 8601 (ex: 2023-10-27T10:30:00.123456789Z ou com offset).")
	}

	dataHoraTransacao = dataHoraTransacao.UTC()
	if dataHoraTransacao.IsZero() {
		logger.AppLogger.Warn("dataHora da transação, apesar de parseável, resultou em uma data/hora zerada (0001-01-01T00:00:00Z).")
		return nil, status.Errorf(codes.InvalidArgument, "dataHora da transação deve ser um valor válido e não zerado")
	}

	agoraUTC := time.Now().UTC()
	if dataHoraTransacao.After(agoraUTC) {
		msgErro := fmt.Sprintf("dataHora da transação (%s) não pode ser no futuro (horário atual do servidor UTC: %s)",
			dataHoraTransacao.Format(time.RFC3339Nano),
			agoraUTC.Format(time.RFC3339Nano))
		logger.AppLogger.Warn(msgErro)
		return nil, status.Errorf(codes.InvalidArgument, "dataHora da transação não pode ser uma data futura")
	}

	s.Armazem.AdicionarTransacao(modelos.Transacoes{
		Valor:         valorTransacao,
		DataTransicao: dataHoraTransacao,
		Uptime:        time.Now().UTC(),
	})

	msgSucesso := fmt.Sprintf("Transação criada com sucesso: dataHora (%s), valor (%.2f)",
		dataHoraTransacao.Format(time.RFC3339Nano),
		valorTransacao)
	logger.AppLogger.Info(msgSucesso)

	return &emptypb.Empty{}, nil
}

func (s *ClientService) GetEstatistica(ctx context.Context, params *server_pb.GetEstatisticaRequest) (res *server_pb.GetEstatisticaResponse, err error) {
	estatisticas := s.Armazem.CalcularEstatisticasUltimoMinuto(params.TimeTravel)
	res = &server_pb.GetEstatisticaResponse{
		Count: int64(estatisticas.Count),
		Sum:   estatisticas.Sum,
		Avg:   estatisticas.Avg,
		Min:   estatisticas.Min,
		Max:   estatisticas.Max,
	}
	return res, nil
}

func (s *ClientService) LimparTransacoes(ctx context.Context, params *emptypb.Empty) (res *emptypb.Empty, err error) {
	s.Armazem.LimparTransacoes()
	return res, nil
}

func (s *ClientService) GetHealthCheck(req *emptypb.Empty, stream server_pb.TransacaoService_GetHealthCheckServer) error {
	for i := range 5 {
		if stream.Context().Err() != nil {
			logger.AppLogger.Warn("Stream cancelado pelo cliente.")
			return stream.Context().Err()
		}

		response := &server_pb.GetHealthCheckResponse{
			Status: fmt.Sprintf("SERVING - Contagem %d", i+1),
		}
		if err := stream.Send(response); err != nil {
			logger.AppLogger.Error("Erro ao enviar stream: %v\n", err)
			return status.Errorf(codes.Internal, "erro ao enviar stream: %v", err)
		}
		logger.AppLogger.Info("Status enviado: %s\n", response.Status)
		time.Sleep(1 * time.Second)
	}

	return nil
}
