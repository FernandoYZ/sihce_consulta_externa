package app

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"sihce_consulta_externa/internal/config"
	"sihce_consulta_externa/internal/config/database"
	page_errors "sihce_consulta_externa/internal/views/pages/errors"
)

// ErrorHandler maneja los errores globales de la aplicación
func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	mensaje := "Ha ocurrido un error interno en el servidor."
	tipo := "INTERNAL_SERVER_ERROR"

	// Si el error es de Fiber, usar su código
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		if code == fiber.StatusNotFound {
			mensaje = "Recurso no encontrado."
			tipo = "NOT_FOUND"
		} else if e.Message != "" {
			mensaje = e.Message
		}
	}

	// Log interno (solo para servidor)
	if code == fiber.StatusInternalServerError {
		log.Printf(
			"[ERROR] %s %s -> %v",
			c.Method(),
			c.Path(),
			err,
		)
	}

	// Si es 404 y la petición espera HTML, renderizar página de error
	acceptsHTML := strings.Contains(c.Get("Accept"), "text/html")
	if code == fiber.StatusNotFound && acceptsHTML {
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.Status(code).SendString(renderNotFoundPage(c))
	}

	// Crear el mapa de respuesta base para JSON
	respuesta := fiber.Map{
		"error":     true,
		"tipo":      tipo,
		"mensaje":   mensaje,
		"path":      c.Path(),
		"metodo":    c.Method(),
		"status":    code,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	// Añadir detalles del error solo en entorno de desarrollo
	cfg := config.ObtenerConfiguracion()
	if cfg != nil && cfg.App.AppEnv == "development" && code == fiber.StatusInternalServerError {
		respuesta["detalles"] = err.Error()
	}

	// Respuesta estandarizada JSON
	return c.Status(code).JSON(respuesta)
}

// renderNotFoundPage renderiza la página 404 usando Templ
func renderNotFoundPage(c *fiber.Ctx) string {
	var buf strings.Builder
	if err := page_errors.NotFound().Render(c.Context(), &buf); err != nil {
		return "<h1>404 - Página no encontrada</h1>"
	}
	return buf.String()
}

// HealthCheckHandler crea un handler para verificar el estado de la API
func HealthCheckHandler(db *database.GestorDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Verificar salud de la BD
		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
		defer cancel()

		dbStatus := "ok"
		if err := db.VerificarSalud(ctx); err != nil {
			log.Printf("Advertencia en health check: %v", err)
			dbStatus = "error"
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": true,
			"data": fiber.Map{
				"mensaje":   "API funcionando correctamente",
				"version":   "1.0.0",
				"timestamp": time.Now().UTC().Format(time.RFC3339),
				"servicios": fiber.Map{
					"database": dbStatus,
				},
			},
		})
	}
}
