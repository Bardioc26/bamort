package main

import (
	"bamort/appsystem"
	"bamort/character"
	"bamort/config"
	"bamort/database"
	"bamort/equipment"
	"bamort/gsmaster"
	"bamort/importer"
	"bamort/importero"
	"bamort/logger"
	"bamort/maintenance"
	"bamort/pdfrender"
	"bamort/router"
	"bamort/transfero"
	"bamort/user"

	"github.com/gin-gonic/gin"

	// Swagger documentation
	_ "bamort/docs/swagger"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title BaMoRT API
// @version 2.0
// @description BaMoRT (MOAM Replacement) - Role-playing Game Character Management System API
// @description
// @description This API provides comprehensive character management, import/export functionality,
// @description master data management, and PDF generation for tabletop role-playing games.
// @description
// @description ## Authentication
// @description Most endpoints require authentication via JWT token. Obtain a token by calling POST /login.
// @description Include the token in the Authorization header: `Authorization: Bearer <token>`
// @description
// @description ## Character Import/Export
// @description The API supports pluggable adapters for importing/exporting characters from various systems.
// @description Supported formats include MOAM VTT, Foundry VTT, and more via microservice adapters.
//
// @contact.name BaMoRT Support
// @contact.email bamort.support@trokan.de
//
// @license.name Custom License
// @license.url https://github.com/Bardioc26/bamort/blob/main/LICENSE
//
// @host localhost:8180
// @BasePath /
// @schemes http https
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token
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

	logger.Info("BaMoRT Server wird gestartet...")
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

	// Run database migrations
	logger.Debug("Führe Datenbank-Migrationen aus...")
	if err := database.MigrateStructure(); err != nil {
		logger.Error("Fehler bei Datenbank-Migrationen: %s", err.Error())
	}
	if err := importer.MigrateStructure(database.DB); err != nil {
		logger.Error("Fehler bei Importer-Migrationen: %s", err.Error())
	} else {
		logger.Info("Datenbank-Migrationen erfolgreich")
	}

	/*
		// Populate initial misc lookup data
		logger.Debug("Initialisiere Misc-Lookup-Daten...")
		if err := gsmaster.PopulateMiscLookupData(); err != nil {
			logger.Warn("Fehler beim Initialisieren der Misc-Lookup-Daten: %s", err.Error())
		} else {
			logger.Info("Misc-Lookup-Daten erfolgreich initialisiert")
		}
	*/

	// Initialize PDF templates
	logger.Debug("Initialisiere PDF-Templates...")
	if err := pdfrender.InitializeTemplates("/app/default_templates", cfg.TemplatesDir); err != nil {
		logger.Warn("Fehler beim Initialisieren der Templates: %s", err.Error())
	} else {
		logger.Info("PDF-Templates erfolgreich initialisiert")
	}

	// Initialize import/export adapter registry
	logger.Debug("Initialisiere Adapter-Registry...")
	adapterRegistry := importer.NewAdapterRegistry()
	importer.InitializeRegistry(adapterRegistry)

	// Register adapters from config (if any)
	// TODO: Load adapters from environment variable IMPORT_ADAPTERS
	// For now, registry is empty and adapters can be registered manually

	// Start background health checker (runs every 30s)
	adapterRegistry.StartBackgroundHealthChecker()
	logger.Info("Adapter-Registry erfolgreich initialisiert und Health-Checker gestartet")

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
	importero.RegisterRoutes(protected)
	importer.RegisterRoutes(protected) // New pluggable import/export system
	pdfrender.RegisterRoutes(protected)
	transfero.RegisterRoutes(protected)
	appsystem.RegisterRoutes(protected)

	// Register public routes (no authentication)
	pdfrender.RegisterPublicRoutes(r)
	appsystem.RegisterPublicRoutes(r)

	// Swagger documentation endpoint
	serverAddress := cfg.GetServerAddress()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	logger.Info("Swagger UI verfügbar unter: %s/swagger/index.html", serverAddress)

	logger.Info("API-Routen erfolgreich registriert")

	// Server starten
	logger.Info("Server startet auf Adresse: %s", serverAddress)
	if err := r.Run(serverAddress); err != nil {
		logger.Error("Fehler beim Starten des Servers: %s", err.Error())
		panic(err)
	}
}
