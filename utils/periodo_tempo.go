package utils

import (
	"desafio-itau-back-grpc/logger"
	"time"
)

func DeterminarJanelaDeTempo(periodo int64) (limiteInferior time.Time) {
	if periodo == 0 {
		periodo = 60
	}
	agora := time.Now()
	offsetEmNanosegundos := -periodo * int64(time.Second)
	offsetComoDuration := time.Duration(offsetEmNanosegundos)
	limiteInferior = agora.Add(offsetComoDuration)

	logger.AppLogger.Info(
		"Janela de tempo determinada: Início=%s, Fim (Agora)=%s, Offset Nanos=%d, Offset Duração=%s",
		limiteInferior.Format(time.RFC3339Nano),
		agora.Format(time.RFC3339Nano),
		offsetEmNanosegundos,
		offsetComoDuration,
	)
	return limiteInferior
}
