package logger

import (
	"fmt"
	"os"
)

func (l *Logger) logf(level Level, format string, v ...interface{}) {
	// calldepth = 2: pula logf -> pula Debugf/Infof/... -> chega no chamador original.
	l.output(level, 2, fmt.Sprintf(format, v...))
}

// Debug formata e loga uma mensagem no nível DEBUG.
func (l *Logger) Debug(format string, v ...interface{}) {
	l.logf(DEBUG, format, v...)
}

// Info formata e loga uma mensagem no nível INFO.
func (l *Logger) Info(format string, v ...interface{}) {
	l.logf(INFO, format, v...)
}

// Warn formata e loga uma mensagem no nível WARN.
func (l *Logger) Warn(format string, v ...interface{}) {
	l.logf(WARN, format, v...)
}

// Error formata e loga uma mensagem no nível ERROR.
func (l *Logger) Error(format string, v ...interface{}) {
	l.logf(ERROR, format, v...)
}

// Fatal formata, loga uma mensagem no nível ERROR e então chama os.Exit(1).
func (l *Logger) Fatal(format string, v ...interface{}) {
	l.logf(ERROR, format, v...)
	l.Close() // Tenta fechar o arquivo de log antes de sair.
	os.Exit(1)
}
