package disk

import "desafio-itau-back-grpc/modelos"

func (a *Armazem) Transacoes() []modelos.Transacoes {
	a.mu.RLock() // Bloqueio para leitura
	defer a.mu.RUnlock()

	copiaClientes := make([]modelos.Transacoes, len(a.transacoes))
	copy(copiaClientes, a.transacoes)
	return copiaClientes
}

func (a *Armazem) AdicionarTransacao(cliente modelos.Transacoes) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.transacoes = append(a.transacoes, cliente)
}
