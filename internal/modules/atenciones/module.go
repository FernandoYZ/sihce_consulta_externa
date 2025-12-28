package atenciones

import (
	"github.com/gofiber/fiber/v2"

	"sihce_consulta_externa/internal/config/database"
	"sihce_consulta_externa/internal/model"
	"sihce_consulta_externa/internal/modules/atenciones/domain"
	repository_atenciones "sihce_consulta_externa/internal/modules/atenciones/repository"
	service_atenciones "sihce_consulta_externa/internal/modules/atenciones/service"
)

// Modulo representa el módulo de atenciones
type Modulo struct {
	controller *atenciones.Controller
}

// NuevoModulo crea una nueva instancia del módulo con inyección de dependencias
func NuevoModulo(db *database.GestorDB, cfg *model.AppConfig) *Modulo {
	// Inyección de dependencias siguiendo el patrón:
	// Repository -> Service -> Controller

	// 1. Crear repository
	repo := repository_atenciones.NuevoRepository(db)

	// 2. Crear service inyectando repository
	service := service_atenciones.NuevoService(repo)

	// 3. Crear controller inyectando service
	controller := atenciones.NuevoController(service)

	return &Modulo{
		controller: controller,
	}
}

// RegistrarRutas registra las rutas del módulo
func (m *Modulo) RegistrarRutas(app *fiber.App, apiRouter fiber.Router) {
	m.controller.RegistrarRutas(app, apiRouter)
}
