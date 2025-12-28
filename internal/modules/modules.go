package modules

import (
	"github.com/gofiber/fiber/v2"

	"sihce_consulta_externa/internal/config/database"
	"sihce_consulta_externa/internal/model"
	"sihce_consulta_externa/internal/modules/atenciones"
)

// Modulos registra todos los módulos de la aplicación
func Modulos(app *fiber.App, apiRouter fiber.Router, db *database.GestorDB, cfg *model.AppConfig) {
	// Registrar módulo de atenciones
	atenciones.NuevoModulo(db, cfg).RegistrarRutas(app, apiRouter)

	// Aquí puedes registrar más módulos en el futuro
	// Ejemplo:
	// pacientes.NuevoModulo(db, cfg).RegistrarRutas(app, apiRouter)
	// citas.NuevoModulo(db, cfg).RegistrarRutas(app, apiRouter)
}
