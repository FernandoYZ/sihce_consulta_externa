package app

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"sihce_consulta_externa/internal/model"
)

// SetupProfiling configura el profiling con pprof
// Solo se habilita en modo desarrollo por razones de seguridad
// IMPORTANTE: Debe ser llamado ANTES de configurar las rutas de la aplicación
func SetupProfiling(app *fiber.App, cfg *model.AppConfig) {
	// Solo habilitar pprof en desarrollo (dev o development)
	isDev := cfg.App.AppEnv == "development" || cfg.App.AppEnv == "dev"
	if !isDev {
		log.Println("⚠️  Profiling deshabilitado en producción")
		return
	}

	// Registrar pprof con configuración default
	// El middleware pprof de Fiber maneja automáticamente todas las rutas
	app.Use(pprof.New())

	log.Println("🔍 Profiling habilitado en /debug/pprof")
	log.Println("📊 Endpoints disponibles:")
	log.Println("   - Index:           /debug/pprof/")
	log.Println("   - CPU Profile:     /debug/pprof/profile?seconds=30")
	log.Println("   - Heap Profile:    /debug/pprof/heap")
	log.Println("   - Goroutines:      /debug/pprof/goroutine")
	log.Println("   - Allocs:          /debug/pprof/allocs")
	log.Println("   - Block:           /debug/pprof/block")
	log.Println("   - Mutex:           /debug/pprof/mutex")
	log.Println("   - Trace:           /debug/pprof/trace?seconds=5")
}
