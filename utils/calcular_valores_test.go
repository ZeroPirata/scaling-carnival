package utils

import (
	"desafio-itau-back-grpc/modelos"
	"math"
	"testing"
)

const floatEqualityThreshold = 1e-9

func insertTransacoes() []modelos.Transacoes {
	transacoes := []modelos.Transacoes{}

	transacoes = append(transacoes, modelos.Transacoes{Valor: 10})
	transacoes = append(transacoes, modelos.Transacoes{Valor: 20})
	transacoes = append(transacoes, modelos.Transacoes{Valor: 30})
	transacoes = append(transacoes, modelos.Transacoes{Valor: 40})
	transacoes = append(transacoes, modelos.Transacoes{Valor: 50})
	transacoes = append(transacoes, modelos.Transacoes{Valor: 60})
	transacoes = append(transacoes, modelos.Transacoes{Valor: 70})
	transacoes = append(transacoes, modelos.Transacoes{Valor: 80})
	transacoes = append(transacoes, modelos.Transacoes{Valor: 90})
	transacoes = append(transacoes, modelos.Transacoes{Valor: 100})

	return transacoes
}

func aproximadamenteEquals(a, b float64) bool {
	return math.Abs(a-b) <= floatEqualityThreshold
}

func TestFuncaoCalcularValores(t *testing.T) {
	// Obter os dados de teste da sua função helper
	transacoesDeTeste := insertTransacoes()

	// Calcular os valores esperados com base em `transacoesDeTeste`
	// Soma: 10+20+30+40+50+60+70+80+90+100 = 550
	// Contagem: 10
	// Média: 550 / 10 = 55
	// Min: 10
	// Max: 100
	esperadoParaSucesso := modelos.EstatisticasTransacao{
		Sum:   550.0,
		Avg:   55.0,
		Min:   10.0,
		Max:   100.0,
		Count: 10,
	}

	testCases := []struct {
		nome       string
		transacoes []modelos.Transacoes // Entrada para a função
		parametro  int64                // Se o parâmetro for realmente usado
		esperado   modelos.EstatisticasTransacao
		esperaErro bool
		msgErro    string
	}{
		{
			nome:       "1 - Sucesso padrão com dados de insertTransacoes",
			transacoes: transacoesDeTeste,
			parametro:  0, // Usado se sua função CalcularValores o utilizar
			esperado:   esperadoParaSucesso,
			esperaErro: false,
		},
		{
			nome:       "2 - Lista de transações vazia",
			transacoes: []modelos.Transacoes{},
			parametro:  0,
			esperado:   modelos.EstatisticasTransacao{}, // ou como sua função trata isso
			esperaErro: true,
			msgErro:    "a lista de transações não pode estar vazia",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.nome, func(t *testing.T) {
			resultado := CalcularEstatisticasParaValores(tc.transacoes)
			if resultado.Count != tc.esperado.Count {
				t.Errorf("Contagem: esperado %d, obtido %d", tc.esperado.Count, resultado.Count)
			}
			if !aproximadamenteEquals(resultado.Sum, tc.esperado.Sum) {
				t.Errorf("Soma: esperado %f, obtido %f", tc.esperado.Sum, resultado.Sum)
			}
			if !aproximadamenteEquals(resultado.Avg, tc.esperado.Avg) {
				t.Errorf("Avg: esperado %f, obtido %f", tc.esperado.Avg, resultado.Avg)
			}
			if !aproximadamenteEquals(resultado.Min, tc.esperado.Min) {
				t.Errorf("Min: esperado %f, obtido %f", tc.esperado.Min, resultado.Min)
			}
			if !aproximadamenteEquals(resultado.Max, tc.esperado.Max) {
				t.Errorf("Max: esperado %f, obtido %f", tc.esperado.Max, resultado.Max)
			}

		})
	}
}
