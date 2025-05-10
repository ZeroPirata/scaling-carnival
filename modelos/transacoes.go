package modelos

import "time"

type Transacoes struct {
	Valor         float64
	DataTransicao time.Time
}

type EstatisticasTransacao struct {
	Count int     `json:"count"`
	Sum   float64 `json:"sum"`
	Avg   float64 `json:"avg"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
}

type TransacaoComTimestamp interface {
	ObterTimestamp() time.Time
}

func (t Transacoes) ObterValorNumerico() float64 {
	return t.Valor
}

func (t Transacoes) ObterTimestamp() time.Time {
	return t.DataTransicao
}
