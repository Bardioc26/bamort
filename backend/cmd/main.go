package main

import (
	"bamort/character"
	"bamort/config"
	"bamort/database"
	"bamort/equipment"
	"bamort/gsmaster"
	"bamort/importer"
	"bamort/logger"
	"bamort/maintenance"
	"bamort/pdfrender"
	"bamort/router"
	"bamort/user"

	"github.com/gin-gonic/gin"
)

// @title Bamort API
// @version 1
// @description This is the API for Bamort
// @host localhost:8180
// @BasePath /
// @schemes http
func main() {
	// Verwende die globale Konfigurationsvariable (bereits in config.init() geladen)
	cfg := config.Cfg

	// Logger konfigurieren
	logger.SetDebugMode(cfg.DebugMode)
	if cfg.LogLevel == "DEBUG" {
		logger.SetMinLogLevel(logger.DEBUG)
	} else if cfg.LogLevel == "WARN" {
		logger.SetMinLogLevel(logger.WARN)
	} else if cfg.LogLevel == "ERROR" {
		logger.SetMinLogLevel(logger.ERROR)
	} else {
		logger.SetMinLogLevel(logger.INFO)
	}

	logger.Info("Bamort Server wird gestartet...")
	logger.Debug("Debug-Modus ist aktiviert")
	logger.Info("Environment: %s", cfg.Environment)
	logger.Info("testingDB Set: %s", cfg.DevTesting)
	logger.Info("Server Port: %s", cfg.ServerPort)

	// Gin-Modus basierend auf Environment setzen
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
		logger.Info("Gin läuft im Release-Modus")
	} else {
		gin.SetMode(gin.DebugMode)
		logger.Debug("Gin läuft im Debug-Modus")
	}

	// Datenbank verbinden
	logger.Debug("Verbinde mit Datenbank...")
	database.ConnectDatabase()
	logger.Info("Datenbankverbindung erfolgreich")

	// Initialize PDF templates
	logger.Debug("Initialisiere PDF-Templates...")
	if err := pdfrender.InitializeTemplates("/app/default_templates", cfg.TemplatesDir); err != nil {
		logger.Warn("Fehler beim Initialisieren der Templates: %s", err.Error())
	} else {
		logger.Info("PDF-Templates erfolgreich initialisiert")
	}

	r := gin.Default()
	router.SetupGin(r)

	// Routes registrieren
	logger.Debug("Registriere API-Routen...")
	protected := router.BaseRouterGrp(r)
	// Register your module routes
	user.RegisterRoutes(protected)
	gsmaster.RegisterRoutes(protected)
	character.RegisterRoutes(protected)
	equipment.RegisterRoutes(protected)
	maintenance.RegisterRoutes(protected)
	importer.RegisterRoutes(protected)
	pdfrender.RegisterRoutes(protected)

	// Register public routes (no authentication)
	pdfrender.RegisterPublicRoutes(r)

	logger.Info("API-Routen erfolgreich registriert")

	// Server starten
	serverAddress := cfg.GetServerAddress()
	logger.Info("Server startet auf Adresse: %s", serverAddress)
	if err := r.Run(serverAddress); err != nil {
		logger.Error("Fehler beim Starten des Servers: %s", err.Error())
		panic(err)
	}
}
