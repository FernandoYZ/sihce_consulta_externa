package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"

	"sihce_consulta_externa/internal/config"
	"sihce_consulta_externa/internal/config/database"
	"sihce_consulta_externa/internal/model"
)

// App encapsula la aplicación con sus dependencias
type App struct {
	Fiber  *fiber.App
	Config *model.AppConfig
	DB     *database.GestorDB
}

// New crea una nueva instancia de la aplicación
func New(cfg *model.AppConfig, db *database.GestorDB) *App {
	// Crear aplicación Fiber con configuración personalizada
	app := fiber.New(fiber.Config{
		AppName:      "SIHCE Consulta Externa v1.0",
		ErrorHandler: ErrorHandler,
	})

	// Configurar middlewares
	Middlewares(app)

	// SetupProfiling(app, cfg)

	// Configurar rutas
	Rutas(app, db, cfg)

	// Configurar archivos estáticos
	ArchivosPublicos(app)

	return &App{
		Fiber:  app,
		Config: cfg,
		DB:     db,
	}
}

// Run inicia la aplicación
func Run() error {
	// Cargar configuración
	if err := config.CargarConfiguracion(); err != nil {
		return fmt.Errorf("error cargando configuración: %w", err)
	}

	cfg := config.ObtenerConfiguracion()

	// Inicializar gestor de base de datos
	gestorDB, err := database.NewGestorDB(&cfg.Database)
	if err != nil {
		return fmt.Errorf("error inicializando gestor de base de datos: %w", err)
	}
	defer gestorDB.ApagarDatabase()

	// Conectar a la base de datos principal ANTES de iniciar el servidor
	if _, err := gestorDB.GetPrincipal(); err != nil {
		return fmt.Errorf("error conectando a la base de datos principal: %w", err)
	}

	// Crear aplicación
	aplicacion := New(cfg, gestorDB)

	// Canal para manejar señales del sistema
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Iniciar servidor en goroutine
	go func() {
		addr := fmt.Sprintf(":%d", cfg.App.Port)
		if err := aplicacion.Fiber.Listen(addr); err != nil {
			log.Fatalf("Error al iniciar servidor: %v", err)
		}
	}()

	// Esperar señal de apagado
	<-quit

	// Cerrar servidor
	return aplicacion.ApagarApp()
}

// Shutdown apaga la aplicación de forma ordenada
func (a *App) ApagarApp() error {
	return a.Fiber.Shutdown()
}
