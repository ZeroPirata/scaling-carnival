package modelos

import (
	"desafio-itau-back-grpc/helper"
	"fmt"
	"io"
	"sync"
)

type Gauge struct {
	nome   string
	ajuda  string
	valor  float64
	mu     sync.RWMutex
	labels map[string]string
}

func NovoGauge(nome, ajuda string, labels map[string]string) *Gauge {
	return &Gauge{nome: nome, ajuda: ajuda, labels: helper.CopiarLabels(labels)}
}
func (g *Gauge) Set(v float64) { g.mu.Lock(); defer g.mu.Unlock(); g.valor = v }
func (g *Gauge) Inc()          { g.mu.Lock(); defer g.mu.Unlock(); g.valor++ }
func (g *Gauge) Dec()          { g.mu.Lock(); defer g.mu.Unlock(); g.valor-- }
func (g *Gauge) Add(v float64) { g.mu.Lock(); defer g.mu.Unlock(); g.valor += v }
func (g *Gauge) Nome() string  { return g.nome }
func (g *Gauge) Ajuda() string { return g.ajuda }
func (g *Gauge) Tipo() string  { return "gauge" }
func (g *Gauge) EscreverExposicao(w io.Writer) error {
	g.mu.RLock()
	defer g.mu.RUnlock()
	labelStr := helper.FormatarLabels(g.labels)
	_, err := fmt.Fprintf(w, "%s%s %f\n", g.nome, labelStr, g.valor)
	return err
}
