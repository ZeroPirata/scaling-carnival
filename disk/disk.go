// Array de Transições
package disk

import (
	"desafio-itau-back-grpc/logger"
	"desafio-itau-back-grpc/modelos"
	"desafio-itau-back-grpc/utils"
	"sync"
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

func (a *Armazem) CalcularEstatisticasUltimoMinuto(limite_tempo int64) modelos.EstatisticasTransacao {
	a.mu.Lock()
	defer a.mu.Unlock()
	periodo := utils.DeterminarJanelaDeTempo(limite_tempo)
	valores := utils.FiltrarElementosPorTempo(a.transacoes, periodo)
	res := utils.CalcularEstatisticasParaValores(valores)
	return res
}

func (a *Armazem) LimparTransacoes() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.transacoes = a.transacoes[:0]
}
