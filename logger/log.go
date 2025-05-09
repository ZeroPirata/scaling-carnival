package logger

import (
	"log"
	"os"
	"path/filepath"
	"sync"
)

var AppLogger *Logger

// Logger define a estrutura do nosso logger personalizado
type Logger struct {
	mu          sync.Mutex  // Para escrita atômica
	fileLogger  *log.Logger // Logger para escrever em arquivo
	logFile     *os.File    // Referência ao arquivo de log
	serviceName string
	debugMode   bool
	minLevel    Level // Nível mínimo para logar
}

func newLogger(serviceName string, defaultMinLevel Level) *Logger {
	return &Logger{
		serviceName: serviceName,
		debugMode:   false, // Será definido em SetupLogging
		minLevel:    defaultMinLevel,
	}
}

func SetupLogging(debug bool, serviceName string, minLevel Level) {
	AppLogger = newLogger(serviceName, minLevel)
	AppLogger.debugMode = debug

	stdLogFlags := log.Ldate | log.Ltime | log.Lshortfile

	if debug {
		// A chamada para tryEnableVirtualTerminalProcessing() foi REMOVIDA.

		log.SetOutput(os.Stdout)
		log.SetFlags(stdLogFlags)
		// AVISO: As cores ANSI podem não funcionar no CMD tradicional do Windows.
		// Funcionarão em terminais como Windows Terminal, Git Bash, Linux, macOS, etc.
		log.Println("Modo Debug HABILITADO: Logs do AppLogger usarão códigos ANSI no console.")
		AppLogger.Info("AppLogger iniciado em modo DEBUG (console com códigos ANSI).")
		return
	}

	// Modo não-debug: Logar em arquivo (lógica igual à anterior).
	logFilePath := os.Getenv("LOG_PATH")
	if logFilePath == "" {
		baseDir := "var/log"
		if _, err := os.Stat(baseDir); os.IsNotExist(err) {
			if mkErr := os.MkdirAll(baseDir, 0755); mkErr != nil {
				log.Printf("ALERTA: Não foi possível criar o diretório base de log '%s': %v", baseDir, mkErr)
			}
		}
		logFilePath = filepath.Join(baseDir, "sql-"+serviceName+".log")
	}

	logDir := filepath.Dir(logFilePath)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			log.SetOutput(os.Stderr)
			log.SetFlags(stdLogFlags)
			log.Printf("ALERTA: Não foi possível criar o diretório de log '%s': %v", logDir, err)
			log.Println("ALERTA: Logs (AppLogger e padrão) continuarão sendo exibidos na Saída de Erro Padrão (Stderr).")
			AppLogger.debugMode = true
			AppLogger.minLevel = DEBUG
			AppLogger.Warn("AppLogger fallback para Stderr devido à falha na criação do diretório de log.")
			return
		}
	}

	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.SetOutput(os.Stderr)
		log.SetFlags(stdLogFlags)
		log.Printf("ALERTA: Não foi possível abrir/criar o arquivo de log '%s': %v", logFilePath, err)
		log.Println("ALERTA: Logs (AppLogger e padrão) continuarão sendo exibidos na Saída de Erro Padrão (Stderr).")
		AppLogger.debugMode = true
		AppLogger.minLevel = DEBUG
		AppLogger.Warn("AppLogger fallback para Stderr devido à falha na abertura do arquivo de log.")
		return
	}

	AppLogger.logFile = logFile
	AppLogger.fileLogger = log.New(logFile, "", stdLogFlags)

	log.SetOutput(logFile)
	log.SetFlags(stdLogFlags)

	log.Println("--- Logging padrão (pacote 'log') iniciado no arquivo ---")
	AppLogger.fileLogger.Printf("[%s] --- AppLogger iniciado no arquivo ---", INFO.String())
}
