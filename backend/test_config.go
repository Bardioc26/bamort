package main

import (
	"bamort/config"
	"fmt"
)

func main() {
	// Teste, ob die globale Konfigurationsvariable funktioniert
	fmt.Printf("Globale Konfiguration:\n")
	fmt.Printf("Environment: %s\n", config.Cfg.Environment)
	fmt.Printf("DatabaseType: %s\n", config.Cfg.DatabaseType)
	fmt.Printf("DatabaseURL: %s\n", config.Cfg.DatabaseURL)
	fmt.Printf("ServerPort: %s\n", config.Cfg.ServerPort)
	fmt.Printf("DebugMode: %v\n", config.Cfg.DebugMode)
	fmt.Printf("LogLevel: %s\n", config.Cfg.LogLevel)
	fmt.Printf("Testing: %s\n", config.Cfg.DevTesting)
}
