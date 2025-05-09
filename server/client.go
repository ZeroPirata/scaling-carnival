package server

import (
	"context"
	"desafio-itau-back-grpc/disk"
	"desafio-itau-back-grpc/logger"
	"desafio-itau-back-grpc/modelos"
	server_pb "desafio-itau-back-grpc/server/pb"
	"fmt"
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

func (s *ClientService) CriarTransacao(ctx context.Context, params *server_pb.CriarTransacaoRequest) (res *emptypb.Empty, err error) {
	valorTransacao := params.GetTransacao().GetValor()
	dataTransacao := params.GetTransacao().GetDataHora()
	if valorTransacao <= 0 {
		msgErro := fmt.Sprintf("valor da transação deve ser positivo, recebido: %.2f", valorTransacao)
		logger.AppLogger.Warn(msgErro)
		return nil, status.Errorf(codes.InvalidArgument, "transação invalida.")
	}
	if dataTransacao == nil {
		msgErro := fmt.Sprintf("valor da data deve ser preenchida: %s", dataTransacao)
		logger.AppLogger.Warn(msgErro)
		return nil, status.Errorf(codes.InvalidArgument, "Data deve ser preenchidas")
	}
	if !dataTransacao.IsValid() {
		msgErro := fmt.Sprintf("valor da data deve ser preenchida de forma correto: %s", dataTransacao)
		logger.AppLogger.Warn(msgErro)
		return nil, status.Errorf(codes.InvalidArgument, "Data deve ser preenchidas com valores validos")
	}
	transacao := dataTransacao.AsTime()
	if transacao.IsZero() {
		msgErro := fmt.Sprintf("valor da data deve ser preenchida de forma correto: %s", dataTransacao)
		logger.AppLogger.Warn(msgErro)
		return nil, status.Errorf(codes.InvalidArgument, "Data deve ser preenchidas com valores validos")
	}
	agora := time.Now()
	if transacao.After(agora) {
		msgErro := fmt.Sprintf("dataHora da transação (%s) não pode ser no futuro (horário atual do servidor: %s)",
			transacao.Format(time.RFC3339),
			agora.Format(time.RFC3339))
		logger.AppLogger.Warn(msgErro)
		return nil, status.Errorf(codes.InvalidArgument, "Data deve ser preenchidas com valores validos")
	}
	s.Armazem.AdicionarTransacao(modelos.Transacoes{
		Valor:         valorTransacao,
		DataTransicao: transacao,
		Uptime:        time.Now(),
	})
	msg := fmt.Sprintf("data hora da transação (%s) | valor (%.2f) ",
		transacao.Format(time.RFC3339),
		valorTransacao)
	logger.AppLogger.Info(msg)
	return res, err
}

func (s *ClientService) GetEstatistica(ctx context.Context, params *emptypb.Empty) (res *server_pb.EstatisticaResponse, err error) {
	estatisticas := s.Armazem.CalcularEstatisticasUltimoMinuto()
	res = &server_pb.EstatisticaResponse{
		Count: int64(estatisticas.Count),
		Sum:   estatisticas.Sum,
		Avg:   estatisticas.Avg,
		Min:   estatisticas.Min,
		Max:   estatisticas.Max,
	}
	return res, nil
}
