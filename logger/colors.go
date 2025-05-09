package logger

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

const (
	AnsiReset  = "\x1b[0m"  // Reseta a formatação
	AnsiRed    = "\x1b[31m" // Vermelho
	AnsiYellow = "\x1b[33m" // Amarelo
	AnsiBlue   = "\x1b[34m" // Azul
	AnsiWhite  = "\x1b[37m" // Branco
)

func (l Level) AnsiColor() string {
	switch l {
	case DEBUG:
		return AnsiBlue
	case INFO:
		return AnsiWhite
	case WARN:
		return AnsiYellow
	case ERROR:
		return AnsiRed
	default:
		return AnsiWhite
	}
}
