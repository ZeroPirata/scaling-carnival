package modelos

import (
	"desafio-itau-back-grpc/helper"
	"fmt"
	"sync"
)

type CounterVec struct {
	nomeBase  string
	ajuda     string
	labelKeys []string            // Ordem definida das chaves de label
	counters  map[string]*Counter // Chave é a string de labels formatada
	mu        sync.RWMutex
	registro  *Registro // Referência ao registro global
}

func NovoCounterVec(reg *Registro, nomeBase, ajuda string, labelKeys []string) *CounterVec {
	return &CounterVec{
		nomeBase:  nomeBase,
		ajuda:     ajuda,
		labelKeys: labelKeys,
		counters:  make(map[string]*Counter),
		registro:  reg, // Importante: para que os counters filhos sejam registrados
	}
}

func (cv *CounterVec) WithLabelValues(valoresLabel ...string) *Counter {
	if len(valoresLabel) != len(cv.labelKeys) {
		panic(fmt.Sprintf("CounterVec %s: esperado %d valores de label, obteve %d", cv.nomeBase, len(cv.labelKeys), len(valoresLabel)))
	}

	labels := make(map[string]string)
	for i, key := range cv.labelKeys {
		labels[key] = valoresLabel[i]
	}
	chaveLabelFormatada := helper.FormatarLabels(labels)

	cv.mu.RLock()
	c, existe := cv.counters[chaveLabelFormatada]
	cv.mu.RUnlock()

	if existe {
		return c
	}

	// Se não existe, cria sob um lock de escrita
	cv.mu.Lock()
	defer cv.mu.Unlock()
	// Verifica novamente, pois outra goroutine pode ter criado enquanto esta esperava o lock
	c, existe = cv.counters[chaveLabelFormatada]
	if existe {
		return c
	}

	novoC := NovoCounter(cv.nomeBase, cv.ajuda, labels)
	cv.counters[chaveLabelFormatada] = novoC
	if cv.registro != nil {
		cv.registro.Registrar(novoC) // Registra o novo counter filho no registro principal
	}
	return novoC
}
