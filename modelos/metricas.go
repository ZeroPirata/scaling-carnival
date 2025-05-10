package modelos

import (
	"io"
)

type Metrica interface {
	Nome() string
	Ajuda() string
	Tipo() string
	EscreverExposicao(writer io.Writer) error
}
