package modelos

import (
	"desafio-itau-back-grpc/helper"
	"fmt"
	"io"
	"sync"
)

type Counter struct {
	nome   string
	ajuda  string
	valor  float64
	mu     sync.RWMutex
	labels map[string]string
}

func NovoCounter(nome, ajuda string, labels map[string]string) *Counter {
	return &Counter{nome: nome, ajuda: ajuda, valor: 0, labels: helper.CopiarLabels(labels)}
}

func (c *Counter) Inc()          { c.mu.Lock(); defer c.mu.Unlock(); c.valor++ }
func (c *Counter) Add(v float64) { c.mu.Lock(); defer c.mu.Unlock(); c.valor += v }
func (c *Counter) Nome() string  { return c.nome }
func (c *Counter) Ajuda() string { return c.ajuda }
func (c *Counter) Tipo() string  { return "counter" }
func (c *Counter) EscreverExposicao(w io.Writer) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	labelStr := helper.FormatarLabels(c.labels)
	_, err := fmt.Fprintf(w, "%s%s %f\n", c.nome, labelStr, c.valor)
	return err
}
