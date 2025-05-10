package modelos

import (
	"fmt"
	"io"
	"sync"
)

type Registro struct {
	metricas []Metrica
	mu       sync.RWMutex
}

func NovoRegistro() *Registro {
	return &Registro{metricas: make([]Metrica, 0)}
}

func (r *Registro) Registrar(m Metrica) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.metricas = append(r.metricas, m)
}

func (r *Registro) ExportarMétricas(writer io.Writer) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ajudasTipos := make(map[string]string)

	for _, m := range r.metricas {
		nomeBase := m.Nome()
		if _, ok := ajudasTipos[nomeBase]; !ok {
			fmt.Fprintf(writer, "# HELP %s %s\n", nomeBase, m.Ajuda())
			fmt.Fprintf(writer, "# TYPE %s %s\n", nomeBase, m.Tipo())
			ajudasTipos[nomeBase] = m.Tipo()
		}
		m.EscreverExposicao(writer)
	}
}
