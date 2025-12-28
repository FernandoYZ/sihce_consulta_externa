package atenciones

import "github.com/gofiber/fiber/v2"

// RegistrarRutas registra las rutas del módulo de atenciones
func (c *Controller) RegistrarRutas(app *fiber.App, apiRouter fiber.Router) {
	// Ruta principal: Página de atenciones (HTML) - Ruta pública
	app.Get("/atenciones", c.RenderizarPaginaAtenciones)  // GET /atenciones

	// Grupo de rutas de API (JSON) - Bajo /api
	api := apiRouter.Group("/atenciones")
	api.Get("/", c.ListarAtenciones)                      // GET /api/atenciones
	api.Get("/estadisticas", c.ObtenerEstadisticas)       // GET /api/atenciones/estadisticas
	api.Get("/:id", c.ObtenerAtencionPorId)               // GET /api/atenciones/:id
}
