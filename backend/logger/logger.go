package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// LogLevel definiert die verschiedenen Log-Level
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// String gibt eine String-Repräsentation des Log-Levels zurück
func (l LogLevel) String() string {
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

// Logger struct enthält die Konfiguration
type Logger struct {
	debugEnabled bool
	minLevel     LogLevel
	logger       *log.Logger
}

// Globaler Logger
var defaultLogger *Logger

// init initialisiert den Standard-Logger
func init() {
	defaultLogger = &Logger{
		debugEnabled: getDebugModeFromEnv(),
		minLevel:     getMinLogLevelFromEnv(),
		logger:       log.New(os.Stdout, "", 0), // Kein Standard-Prefix, wir verwenden unseren eigenen
	}
}

// getDebugModeFromEnv liest den Debug-Modus aus Umgebungsvariablen
func getDebugModeFromEnv() bool {
	debugMode := os.Getenv("DEBUG")
	return strings.ToLower(debugMode) == "true" || debugMode == "1"
}

// getMinLogLevelFromEnv liest das minimale Log-Level aus Umgebungsvariablen
func getMinLogLevelFromEnv() LogLevel {
	levelStr := strings.ToUpper(os.Getenv("LOG_LEVEL"))
	switch levelStr {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN":
		return WARN
	case "ERROR":
		return ERROR
	default:
		if getDebugModeFromEnv() {
			return DEBUG
		}
		return INFO
	}
}

// SetDebugMode aktiviert oder deaktiviert den Debug-Modus
func SetDebugMode(enabled bool) {
	defaultLogger.debugEnabled = enabled
	if enabled && defaultLogger.minLevel > DEBUG {
		defaultLogger.minLevel = DEBUG
	}
}

// SetMinLogLevel setzt das minimale Log-Level
func SetMinLogLevel(level LogLevel) {
	defaultLogger.minLevel = level
}

// IsDebugEnabled gibt zurück, ob der Debug-Modus aktiviert ist
func IsDebugEnabled() bool {
	return defaultLogger.debugEnabled && defaultLogger.minLevel <= DEBUG
}

// logMessage ist die interne Funktion zum Loggen von Nachrichten
func logMessage(level LogLevel, format string, args ...interface{}) {
	if level < defaultLogger.minLevel {
		return
	}

	// Für Debug-Messages zusätzlich prüfen, ob Debug-Modus aktiviert ist
	if level == DEBUG && !defaultLogger.debugEnabled {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("[%s] %s: %s", timestamp, level.String(), message)

	defaultLogger.logger.Println(logLine)
}

// Debug loggt Debug-Nachrichten (nur wenn Debug-Modus aktiviert)
func Debug(format string, args ...interface{}) {
	logMessage(DEBUG, format, args...)
}

// Info loggt Info-Nachrichten
func Info(format string, args ...interface{}) {
	logMessage(INFO, format, args...)
}

// Warn loggt Warn-Nachrichten
func Warn(format string, args ...interface{}) {
	logMessage(WARN, format, args...)
}

// Error loggt Error-Nachrichten
func Error(format string, args ...interface{}) {
	logMessage(ERROR, format, args...)
}

// Debugf ist ein Alias für Debug (für Kompatibilität)
func Debugf(format string, args ...interface{}) {
	Debug(format, args...)
}

// Infof ist ein Alias für Info (für Kompatibilität)
func Infof(format string, args ...interface{}) {
	Info(format, args...)
}

// Warnf ist ein Alias für Warn (für Kompatibilität)
func Warnf(format string, args ...interface{}) {
	Warn(format, args...)
}

// Errorf ist ein Alias für Error (für Kompatibilität)
func Errorf(format string, args ...interface{}) {
	Error(format, args...)
}
