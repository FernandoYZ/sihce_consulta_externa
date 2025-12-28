package app

import (
	"sihce_consulta_externa/internal/config/database"
	"sihce_consulta_externa/internal/model"
	"sihce_consulta_externa/internal/modules"

	"github.com/gofiber/fiber/v2"
)

func Rutas(app *fiber.App, db *database.GestorDB, cfg *model.AppConfig) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("¡Bienvenido a SIHCE Consulta Externa!")
	})

	api := app.Group("/api")

	api.Get("/", HealthCheckHandler(db))

	modules.Modulos(app, api, db, cfg)
}