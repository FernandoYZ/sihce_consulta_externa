package app

import "github.com/gofiber/fiber/v2"

func ArchivosPublicos(app *fiber.App) {
	app.Static("/", "./public")
}