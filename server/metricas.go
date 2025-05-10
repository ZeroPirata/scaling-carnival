package server

import (
	"desafio-itau-back-grpc/modelos"
	"log"
	"net/http"
)

func IniciarServidorMetricas(endereco string, registro *modelos.Registro) {
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		registro.ExportarMétricas(w) // Usa o método do seu Registro
	})

	log.Printf("Servidor de métricas ouvindo em %s", endereco)
	if err := http.ListenAndServe(endereco, mux); err != nil {
		log.Fatalf("Falha ao iniciar servidor de métricas: %v", err)
	}
}
