package utils

import (
	"desafio-itau-back-grpc/modelos"
	"reflect"
	"testing"
	"time"
)

func makeTime(ano int, mes time.Month, dia, hora, min, seg int) time.Time {
	return time.Date(ano, mes, dia, hora, min, seg, 0, time.UTC)
}

func TestFiltrarElementosPorTempo(t *testing.T) {
	limiteBase := makeTime(2025, time.May, 10, 12, 0, 0)

	elemento1Antes := modelos.Transacoes{DataTransicao: limiteBase.Add(-2 * time.Hour)}
	elemento2NoLimite := modelos.Transacoes{DataTransicao: limiteBase}
	elemento3Depois := modelos.Transacoes{DataTransicao: limiteBase.Add(1 * time.Hour)}
	elemento4BemDepois := modelos.Transacoes{DataTransicao: limiteBase.Add(5 * time.Hour)}

	testCases := []struct {
		nome           string
		todosElementos []modelos.Transacoes
		limiteInferior time.Time
		esperado       []modelos.Transacoes // Pode ser nil
	}{
		{
			nome:           "Lista vazia",
			todosElementos: []modelos.Transacoes{},
			limiteInferior: limiteBase,
			esperado:       nil, // CORRIGIDO: Função retorna nil para entrada de slice vazio sem correspondências
		},
		{
			nome:           "Nenhum elemento na lista",
			todosElementos: nil, // Entrada nil
			limiteInferior: limiteBase,
			esperado:       nil, // Função retorna nil para entrada nil
		},
		{
			nome:           "Todos os elementos depois do limite",
			todosElementos: []modelos.Transacoes{elemento3Depois, elemento4BemDepois},
			limiteInferior: limiteBase,
			esperado:       []modelos.Transacoes{elemento3Depois, elemento4BemDepois},
		},
		{
			nome:           "Nenhum elemento depois do limite (todos antes ou no limite)",
			todosElementos: []modelos.Transacoes{elemento1Antes, elemento2NoLimite},
			limiteInferior: limiteBase,
			esperado:       nil, // CORRIGIDO: Função retorna nil se nada for filtrado
		},
		{
			nome: "Alguns elementos depois do limite",
			todosElementos: []modelos.Transacoes{
				elemento1Antes,
				elemento2NoLimite,
				elemento3Depois,
				elemento4BemDepois,
			},
			limiteInferior: limiteBase,
			esperado:       []modelos.Transacoes{elemento3Depois, elemento4BemDepois},
		},
		{
			nome:           "Limite no passado distante (todos os elementos são depois)",
			todosElementos: []modelos.Transacoes{elemento1Antes, elemento2NoLimite, elemento3Depois},
			limiteInferior: makeTime(2000, time.January, 1, 0, 0, 0),
			// Se todos os elementos realmente tiverem DataTransicao após 2000-01-01,
			// então todos devem ser retornados.
			esperado: []modelos.Transacoes{elemento1Antes, elemento2NoLimite, elemento3Depois},
		},
		{
			nome:           "Limite no futuro distante (nenhum elemento é depois)",
			todosElementos: []modelos.Transacoes{elemento1Antes, elemento2NoLimite, elemento3Depois},
			limiteInferior: makeTime(2050, time.January, 1, 0, 0, 0),
			esperado:       nil, // CORRIGIDO: Função retorna nil se nada for filtrado
		},
		{
			nome: "Elementos com o mesmo DataTransicao, alguns depois do limite",
			todosElementos: []modelos.Transacoes{
				{DataTransicao: limiteBase.Add(1 * time.Hour)},  // Depois
				{DataTransicao: limiteBase.Add(-1 * time.Hour)}, // Antes
				{DataTransicao: limiteBase.Add(1 * time.Hour)},  // Depois
			},
			limiteInferior: limiteBase,
			esperado: []modelos.Transacoes{
				{DataTransicao: limiteBase.Add(1 * time.Hour)},
				{DataTransicao: limiteBase.Add(1 * time.Hour)},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.nome, func(t *testing.T) {
			// Assumindo que FiltrarElementosPorTempo está no mesmo pacote 'utils'
			// ou importado corretamente se estiver em outro lugar.
			// Se `FiltrarElementosPorTempo` estiver no pacote `utils` (onde está o teste),
			// a chamada é direta. Se estiver em outro, seria `outro_pacote.FiltrarElementosPorTempo`.
			// Pelo seu comando `desafio-itau-back-grpc/utils`, parece que a função e o teste
			// estão no mesmo pacote `utils`.
			resultado := FiltrarElementosPorTempo(tc.todosElementos, tc.limiteInferior)

			if !reflect.DeepEqual(resultado, tc.esperado) {
				t.Errorf("Teste '%s' falhou:\nLimite Inferior: %v\nEntrada:        %+v\nEsperado:       %+v\nObtido:         %+v",
					tc.nome, tc.limiteInferior.Format(time.RFC3339), tc.todosElementos, tc.esperado, resultado)
			}
		})
	}
}
