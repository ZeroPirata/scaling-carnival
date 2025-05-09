package utils

import (
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ValidarEntrada(campos map[string]any) error {
	for campo, valor := range campos {
		switch v := valor.(type) {
		case string:
			if v == "" {
				return status.Errorf(codes.InvalidArgument, "o campo '%s' não pode ser vazio", campo)
			}
		case time.Time:
			if v.IsZero() {
				return status.Errorf(codes.InvalidArgument, "o campo '%s' não pode ser uma data vazia", campo)
			}
		}

	}
	return nil
}
