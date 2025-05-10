package utils

import (
	"desafio-itau-back-grpc/modelos"
	"time"
)

func FiltrarElementosPorTempo[E modelos.TransacaoComTimestamp](todosElementos []E, limiteInferior time.Time) []E {
	var elementosFiltrados []E
	for _, el := range todosElementos {
		if el.ObterTimestamp().After(limiteInferior) {
			elementosFiltrados = append(elementosFiltrados, el)
		}
	}
	return elementosFiltrados
}
