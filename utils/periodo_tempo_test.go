package utils

import (
	reallogger "desafio-itau-back-grpc/logger"
	"testing"
	"time"

	"bou.ke/monkey"
)

func TestDeterminarJanelaDeTempo(t *testing.T) {

	originalGlobalLogger := reallogger.AppLogger

	if reallogger.AppLogger == nil {
		t.Logf("AVISO: reallogger.AppLogger estava nil. Tentando atribuir uma instância placeholder. Esta é uma solução de contorno e pode não ser suficiente ou correta para o seu logger específico.")

		reallogger.AppLogger = &reallogger.Logger{}

		if reallogger.AppLogger == nil {

			t.Fatalf("Falha crítica: Não foi possível inicializar reallogger.AppLogger. Verifique a API do seu pacote de logger e adapte a configuração do logger neste teste.")
		}
		t.Logf("reallogger.AppLogger foi atribuído a uma instância placeholder: %+v", reallogger.AppLogger)
	}

	t.Cleanup(func() {

		reallogger.AppLogger = originalGlobalLogger
	})

	tempoFixo := time.Date(2025, time.May, 10, 15, 30, 0, 0, time.UTC)
	patchGuard := monkey.Patch(time.Now, func() time.Time {
		return tempoFixo
	})
	t.Cleanup(func() {
		patchGuard.Unpatch()
	})

	testCases := []struct {
		nome                   string
		periodoEntrada         int64
		periodoEfetivoEsperado int64
	}{
		{
			nome:                   "Período zero (deve usar padrão de 60s)",
			periodoEntrada:         0,
			periodoEfetivoEsperado: 60,
		},
		{
			nome:                   "Período positivo (30s)",
			periodoEntrada:         30,
			periodoEfetivoEsperado: 30,
		},
		{
			nome:                   "Período positivo grande (3600s = 1 hora)",
			periodoEntrada:         3600,
			periodoEfetivoEsperado: 3600,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.nome, func(t *testing.T) {

			limiteInferiorObtido := DeterminarJanelaDeTempo(tc.periodoEntrada)

			offsetEsperadoEmNanosegundos := -tc.periodoEfetivoEsperado * int64(time.Second)
			offsetEsperadoComoDuration := time.Duration(offsetEsperadoEmNanosegundos)
			limiteInferiorEsperado := tempoFixo.Add(offsetEsperadoComoDuration)

			if !limiteInferiorObtido.Equal(limiteInferiorEsperado) {
				t.Errorf("Teste '%s' falhou:\nPeríodo de Entrada: %d\nTempo Fixo (Agora Mockado): %s\nEsperado Limite Inferior: %s [Local: %s]\nObtido Limite Inferior:   %s [Local: %s]",
					tc.nome,
					tc.periodoEntrada,
					tempoFixo.Format(time.RFC3339Nano),
					limiteInferiorEsperado.Format(time.RFC3339Nano),
					limiteInferiorEsperado.Local().Format(time.RFC3339Nano),
					limiteInferiorObtido.Format(time.RFC3339Nano),
					limiteInferiorObtido.Local().Format(time.RFC3339Nano),
				)
			}
		})
	}
}
