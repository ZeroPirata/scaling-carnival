package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func (l *Logger) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.logFile != nil {
		if l.fileLogger != nil {
			l.fileLogger.Printf("[%s] --- Logging finalizado no arquivo ---", INFO.String())
		}
		l.logFile.Close()
		l.logFile = nil
		l.fileLogger = nil // Garante que não tentaremos mais logar no arquivo fechado.
	}
}

// output é a função central de logging do AppLogger.
func (l *Logger) output(level Level, calldepth int, msg string) {
	if level < l.minLevel {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now()
	levelStr := level.String()

	if l.debugMode {
		// Saída para o console com códigos ANSI.
		var file string
		var line int
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth + 1)
		if !ok {
			file = "???"
			line = 0
		} else {
			file = filepath.Base(file)
		}

		ansiColorCode := level.AnsiColor()

		consoleFormattedMsg := fmt.Sprintf("%s%s %s:%d [%s] %s%s",
			ansiColorCode,
			timestamp.Format("2006/01/02 15:04:05"),
			file, line,
			levelStr,
			strings.TrimSuffix(msg, "\n"),
			AnsiReset,
		)
		fmt.Println(consoleFormattedMsg)

	} else if l.fileLogger != nil {
		// Saída para arquivo (sem cores).
		err := l.fileLogger.Output(calldepth+1, fmt.Sprintf("[%s] %s", levelStr, msg))
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERRO CRÍTICO AO ESCREVER NO ARQUIVO DE LOG: %v. Mensagem original: [%s] %s\n", err, levelStr, msg)
		}
	}
}
