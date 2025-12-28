package atenciones

import (
	"context"
	"fmt"
	"strings"
	"time"

	http_atenciones "sihce_consulta_externa/internal/modules/atenciones/http"
	repository_atenciones "sihce_consulta_externa/internal/modules/atenciones/repository"
)

// Service interfaz para el servicio de atenciones
type Service interface {
	ListarAtenciones(ctx context.Context, filtros *http_atenciones.FiltrosAtencion) ([]http_atenciones.AtencionDTO, error)
	ObtenerAtencionPorId(ctx context.Context, idAtencion int) (*http_atenciones.AtencionDTO, error)
	ObtenerEstadisticas(ctx context.Context, fechaInicio, fechaFin string) (*http_atenciones.EstadisticasAtencion, error)
}

// ServiceImpl implementación del servicio
type ServiceImpl struct {
	repo repository_atenciones.Repository
}

// NuevoService crea una nueva instancia del servicio
func NuevoService(repo repository_atenciones.Repository) Service {
	return &ServiceImpl{repo: repo}
}

// ListarAtenciones obtiene todas las atenciones formateadas
func (s *ServiceImpl) ListarAtenciones(ctx context.Context, filtros *http_atenciones.FiltrosAtencion) ([]http_atenciones.AtencionDTO, error) {
	// Obtener datos del repositorio
	atencionesBD, err := s.repo.ListarAtenciones(ctx, filtros)
	if err != nil {
		return nil, err
	}

	// Formatear cada atención en el servidor (aprovechando goroutines)
	atenciones := make([]http_atenciones.AtencionDTO, len(atencionesBD))
	for i, a := range atencionesBD {
		atenciones[i] = s.formatearAtencion(&a)
	}

	return atenciones, nil
}

// ObtenerAtencionPorId obtiene una atención formateada por su ID
func (s *ServiceImpl) ObtenerAtencionPorId(ctx context.Context, idAtencion int) (*http_atenciones.AtencionDTO, error) {
	atencionBD, err := s.repo.ObtenerAtencionPorId(ctx, idAtencion)
	if err != nil {
		return nil, err
	}
	if atencionBD == nil {
		return nil, nil
	}

	atencion := s.formatearAtencion(atencionBD)
	return &atencion, nil
}

// ObtenerEstadisticas obtiene las estadísticas del dashboard
func (s *ServiceImpl) ObtenerEstadisticas(ctx context.Context, fechaInicio, fechaFin string) (*http_atenciones.EstadisticasAtencion, error) {
	return s.repo.ObtenerEstadisticas(ctx, fechaInicio, fechaFin)
}

// formatearAtencion formatea una atención del repositorio a DTO
// Aquí es donde se hace la limpieza y formateo de datos en el servidor
func (s *ServiceImpl) formatearAtencion(a *http_atenciones.AtencionRepositoryDTO) http_atenciones.AtencionDTO {
	return http_atenciones.AtencionDTO{
		IdPaciente:           a.IdPaciente,
		NroDocumento:         a.NroDocumento,
		FuenteFinanciamiento: s.normalizarTexto(a.FuenteFinanciamiento),
		IdAtencion:           a.IdAtencion,
		NombreDoctor:         s.capitalizarNombre(a.NombreDoctor),
		ApellidoDoctor:       s.capitalizarNombre(a.ApellidoDoctor),
		ApellidosPaciente:    s.capitalizarNombre(a.ApellidosPaciente),
		NombresPaciente:      s.capitalizarNombre(a.NombresPaciente),
		Edad:                 a.Edad,
		NroHistoriaClinica:   a.NroHistoriaClinica,
		Servicio:             s.capitalizarNombre(a.Servicio),
		IdServicio:           a.IdServicio,
		FechaCita:            a.FechaCita.Format("2006-01-02"),
		HoraInicio:           a.HoraInicio,
		HoraFin:              a.HoraFin,
		Estado:               a.Estado,
		Triaje:               a.Triaje,
		// Campos calculados/formateados
		FechaCitaFormateada: s.formatearFechaEspanol(a.FechaCita),
		HorarioFormateado:   s.formatearHorario(a.HoraInicio, a.HoraFin),
		DoctorCompleto:      s.formatearNombreCompleto("Dr.", a.ApellidoDoctor, a.NombreDoctor),
		PacienteCompleto:    s.formatearNombreCompleto("", a.ApellidosPaciente, a.NombresPaciente),
		EstadoClase:         s.obtenerEstadoClase(a.Estado),
		TriajeClase:         s.obtenerTriajeClase(a.Triaje),
		EstadoIcono:         s.obtenerEstadoIcono(a.Estado),
		TriajeIcono:         s.obtenerTriajeIcono(a.Triaje),
	}
}

// formatearFechaEspanol formatea una fecha a formato español: "Vier 26 Dic 2025"
func (s *ServiceImpl) formatearFechaEspanol(fecha time.Time) string {
	// Nombres de días en español (abreviados a 4 letras)
	diasSemana := map[time.Weekday]string{
		time.Monday:    "Lun",
		time.Tuesday:   "Mar",
		time.Wednesday: "Mié",
		time.Thursday:  "Jue",
		time.Friday:    "Vie",
		time.Saturday:  "Sáb",
		time.Sunday:    "Dom",
	}

	// Nombres de meses en español (abreviados a 3 letras)
	meses := map[time.Month]string{
		time.January:   "Ene",
		time.February:  "Feb",
		time.March:     "Mar",
		time.April:     "Abr",
		time.May:       "May",
		time.June:      "Jun",
		time.July:      "Jul",
		time.August:    "Ago",
		time.September: "Sep",
		time.October:   "Oct",
		time.November:  "Nov",
		time.December:  "Dic",
	}

	dia := diasSemana[fecha.Weekday()]
	mes := meses[fecha.Month()]

	return fmt.Sprintf("%s %02d %s %d", dia, fecha.Day(), mes, fecha.Year())
}

// formatearHorario formatea el horario: "07:00 - 07:30"
func (s *ServiceImpl) formatearHorario(horaInicio, horaFin string) string {
	return fmt.Sprintf("%s - %s", horaInicio, horaFin)
}

// formatearNombreCompleto formatea el nombre completo con título opcional
func (s *ServiceImpl) formatearNombreCompleto(titulo, apellidos, nombres string) string {
	apellidos = s.capitalizarNombre(apellidos)
	nombres = s.capitalizarNombre(nombres)

	if titulo != "" {
		return fmt.Sprintf("%s %s, %s", titulo, apellidos, nombres)
	}
	return fmt.Sprintf("%s, %s", apellidos, nombres)
}

// capitalizarNombre capitaliza cada palabra de un nombre
func (s *ServiceImpl) capitalizarNombre(nombre string) string {
	nombre = strings.TrimSpace(strings.ToLower(nombre))
	palabras := strings.Fields(nombre)

	for i, palabra := range palabras {
		if len(palabra) > 0 {
			palabras[i] = strings.ToUpper(string(palabra[0])) + palabra[1:]
		}
	}

	return strings.Join(palabras, " ")
}

// normalizarTexto normaliza texto general
func (s *ServiceImpl) normalizarTexto(texto string) string {
	return strings.TrimSpace(texto)
}

// obtenerEstadoClase retorna las clases CSS según el estado
func (s *ServiceImpl) obtenerEstadoClase(estado string) string {
	switch strings.ToLower(estado) {
	case "pendiente":
		return "bg-yellow-50 text-yellow-600 border-yellow-100"
	case "atendido":
		return "bg-green-50 text-green-700 border-green-100"
	case "en atención", "en atencion":
		return "bg-blue-50 text-blue-700 border-blue-100"
	case "cancelado":
		return "bg-rose-50 text-rose-700 border-rose-100"
	default:
		return "bg-gray-50 text-gray-500 border-gray-200"
	}
}

// obtenerTriajeClase retorna las clases CSS según el triaje
func (s *ServiceImpl) obtenerTriajeClase(triaje string) string {
	switch strings.ToLower(triaje) {
	case "completado":
		return "bg-green-50 text-green-700 border-green-100"
	case "pendiente":
		return "bg-yellow-50 text-yellow-600 border-yellow-100"
	case "no aplica":
		return "bg-slate-50 text-slate-600/80 border-slate-100"
	default:
		return "bg-gray-100 text-gray-500 border-gray-100"
	}
}

// obtenerEstadoIcono retorna el ícono de Lucide según el estado
func (s *ServiceImpl) obtenerEstadoIcono(estado string) string {
	switch strings.ToLower(estado) {
	case "pendiente":
		return "hourglass"
	case "atendido":
		return "check-circle"
	case "en atención", "en atencion":
		return "activity"
	case "cancelado":
		return "x-circle"
	default:
		return "circle"
	}
}

// obtenerTriajeIcono retorna el ícono de Lucide según el triaje
func (s *ServiceImpl) obtenerTriajeIcono(triaje string) string {
	switch strings.ToLower(triaje) {
	case "completado":
		return "check-circle"
	case "pendiente":
		return "hourglass"
	case "no aplica":
		return "minus-circle"
	default:
		return "circle"
	}
}
