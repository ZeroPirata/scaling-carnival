package bootstrap

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Printf("[WARN] Erro ao carregar .env: %v", err)
	}
}
