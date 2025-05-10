// Array de Transições
package disk

import (
	"desafio-itau-back-grpc/logger"
	"desafio-itau-back-grpc/modelos"
	"sync"
	"time"
)

type Armazem struct {
	transacoes []modelos.Transacoes
	mu         sync.RWMutex
}

var (
	instanciaUnica *Armazem
	once           sync.Once
)

func GetInstanciaUnica() *Armazem {
	logger.AppLogger.Info("Armazenamento interno iniciado.")
	once.Do(func() {
		instanciaUnica = &Armazem{
			transacoes: []modelos.Transacoes{},
		}
	})
	return instanciaUnica
}

func (a *Armazem) CalcularEstatisticasUltimoMinuto() modelos.EstatisticasTransacao {
	a.mu.Lock()
	defer a.mu.Unlock()

	agora := time.Now()
	limiteInferiorTempo := agora.Add(-60 * time.Second)

	var count int = 0
	var sum float64 = 0.0
	var minVal float64 = 0.0
	var maxVal float64 = 0.0
	primeiraTransacaoNoPeriodo := true

	logger.AppLogger.Info("Calculando estatísticas. Agora: %s, Limite de tempo: %s\n", agora.Format(time.RFC3339), limiteInferiorTempo.Format(time.RFC3339))

	for _, t := range a.transacoes {
		if t.DataTransicao.After(limiteInferiorTempo) {
			logger.AppLogger.Info("Transação incluída no cálculo: Valor=%.2f, DataHora=%s\n", t.Valor, t.DataTransicao.Format(time.RFC3339))
			if primeiraTransacaoNoPeriodo {
				minVal = t.Valor
				maxVal = t.Valor
				primeiraTransacaoNoPeriodo = false
			} else {
				if t.Valor < minVal {
					minVal = t.Valor
				}
				if t.Valor > maxVal {
					maxVal = t.Valor
				}
			}
			sum += t.Valor
			count++
		}
	}

	var avg float64 = 0.0
	if count > 0 {
		avg = sum / float64(count)
	}

	return modelos.EstatisticasTransacao{
		Count: count,
		Sum:   sum,
		Avg:   avg,
		Min:   minVal,
		Max:   maxVal,
	}
}

func (a *Armazem) LimparTransacoes() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.transacoes = a.transacoes[:0]
}
