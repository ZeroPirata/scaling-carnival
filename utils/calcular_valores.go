package utils

import "desafio-itau-back-grpc/modelos"

type EntidadeComValor interface {
	ObterValorNumerico() float64
}

func CalcularEstatisticasParaValores[E EntidadeComValor](elementos []E) modelos.EstatisticasTransacao {
	if len(elementos) == 0 {
		return modelos.EstatisticasTransacao{
			Count: 0,
			Sum:   0.0,
			Avg:   0.0,
			Min:   0.0,
			Max:   0.0,
		}
	}

	minVal := elementos[0].ObterValorNumerico()
	maxVal := elementos[0].ObterValorNumerico()

	var count int = 0
	var sum float64 = 0.0

	for _, el := range elementos {
		valorAtual := el.ObterValorNumerico()

		if valorAtual < minVal {
			minVal = valorAtual
		}
		if valorAtual > maxVal {
			maxVal = valorAtual
		}
		sum += valorAtual
		count++
	}

	avg := sum / float64(count)

	return modelos.EstatisticasTransacao{
		Count: count,
		Sum:   sum,
		Avg:   avg,
		Min:   minVal,
		Max:   maxVal,
	}
}
