package atenciones

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	http_atenciones "sihce_consulta_externa/internal/modules/atenciones/http"
	service_atenciones "sihce_consulta_externa/internal/modules/atenciones/service"
	page_atenciones "sihce_consulta_externa/internal/views/pages/atenciones"
)

// Controller maneja las peticiones HTTP para atenciones
type Controller struct {
	service service_atenciones.Service
}

// NuevoController crea una nueva instancia del controller
func NuevoController(service service_atenciones.Service) *Controller {
	return &Controller{service: service}
}

// ListarAtenciones maneja GET /atenciones
func (c *Controller) ListarAtenciones(ctx *fiber.Ctx) error {
	// Parsear filtros desde query params
	filtros := &http_atenciones.FiltrosAtencion{}
	if err := ctx.QueryParser(filtros); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al parsear filtros",
		})
	}

	// Obtener atenciones del servicio
	atenciones, err := c.service.ListarAtenciones(ctx.Context(), filtros)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al obtener atenciones",
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    atenciones,
		"total":   len(atenciones),
	})
}

// ObtenerAtencionPorId maneja GET /atenciones/:id
func (c *Controller) ObtenerAtencionPorId(ctx *fiber.Ctx) error {
	// Parsear ID
	idParam := ctx.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID inválido",
		})
	}

	// Obtener atención del servicio
	atencion, err := c.service.ObtenerAtencionPorId(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al obtener atención",
		})
	}

	if atencion == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Atención no encontrada",
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    atencion,
	})
}

// ObtenerEstadisticas maneja GET /atenciones/estadisticas
func (c *Controller) ObtenerEstadisticas(ctx *fiber.Ctx) error {
	// Obtener fechas de los query params (opcional)
	fechaInicio := ctx.Query("fechaInicio", "")
	fechaFin := ctx.Query("fechaFin", "")

	// Obtener estadísticas del servicio
	stats, err := c.service.ObtenerEstadisticas(ctx.Context(), fechaInicio, fechaFin)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al obtener estadísticas",
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    stats,
	})
}

// RenderizarPaginaAtenciones maneja GET /atenciones (retorna HTML con Templ)
func (c *Controller) RenderizarPaginaAtenciones(ctx *fiber.Ctx) error {
	// Obtener estadísticas
	stats, err := c.service.ObtenerEstadisticas(ctx.Context(), "", "")
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error al cargar estadísticas")
	}

	// Obtener atenciones
	atenciones, err := c.service.ListarAtenciones(ctx.Context(), &http_atenciones.FiltrosAtencion{})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error al cargar atenciones")
	}

	// Renderizar con Templ
	ctx.Set("Content-Type", "text/html; charset=utf-8")

	// Renderizar la página de atenciones usando Templ
	return page_atenciones.Index(stats, atenciones).Render(ctx.Context(), ctx.Response().BodyWriter())
}
