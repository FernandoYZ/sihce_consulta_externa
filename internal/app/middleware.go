package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Middlewares(app *fiber.App) {
	// Recover middleware para recuperarse de panics
	app.Use(recover.New())

	// Logger middleware para logging de peticiones
	app.Use(logger.New())
}