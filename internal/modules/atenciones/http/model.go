package atenciones

import "time"

// AtencionDTO representa una atención en consulta externa
type AtencionDTO struct {
	IdPaciente           int    `json:"idPaciente"`
	NroDocumento         string `json:"nroDocumento"`
	FuenteFinanciamiento string `json:"fuenteFinanciamiento"`
	IdAtencion           int    `json:"idAtencion"`
	NombreDoctor         string `json:"nombreDoctor"`
	ApellidoDoctor       string `json:"apellidoDoctor"`
	ApellidosPaciente    string `json:"apellidosPaciente"`
	NombresPaciente      string `json:"nombresPaciente"`
	Edad                 int    `json:"edad"`
	NroHistoriaClinica   int    `json:"nroHistoriaClinica"`
	Servicio             string `json:"servicio"`
	IdServicio           int    `json:"idServicio"`
	FechaCita            string `json:"fechaCita"`
	HoraInicio           string `json:"horaInicio"`
	HoraFin              string `json:"horaFin"`
	Estado               string `json:"estado"`
	Triaje               string `json:"triaje"`

	// Campos calculados/formateados (procesados en servidor)
	FechaCitaFormateada string `json:"fechaCitaFormateada"` // "Vier 26 Dic 2025"
	HorarioFormateado   string `json:"horarioFormateado"`   // "07:00 - 07:30"
	DoctorCompleto      string `json:"doctorCompleto"`      // "Dr. Apellido, Nombre"
	PacienteCompleto    string `json:"pacienteCompleto"`    // "Apellidos, Nombres"
	EstadoClase         string `json:"estadoClase"`         // Clase CSS para estado
	TriajeClase         string `json:"triajeClase"`         // Clase CSS para triaje
	EstadoIcono         string `json:"estadoIcono"`         // Icono Lucide
	TriajeIcono         string `json:"triajeIcono"`         // Icono Lucide
}

// AtencionRepositoryDTO es el DTO que viene directo de la BD (sin formatear)
type AtencionRepositoryDTO struct {
	IdPaciente           int
	NroDocumento         string
	FuenteFinanciamiento string
	IdAtencion           int
	NombreDoctor         string
	ApellidoDoctor       string
	ApellidosPaciente    string
	NombresPaciente      string
	Edad                 int
	NroHistoriaClinica   int
	Servicio             string
	IdServicio           int
	FechaCita            time.Time
	HoraInicio           string
	HoraFin              string
	Estado               string
	Triaje               string
}

// FiltrosAtencion representa los filtros para búsqueda de atenciones
type FiltrosAtencion struct {
	FechaInicio          string `query:"fechaInicio"`
	FechaFin             string `query:"fechaFin"`
	NroHistoriaClinica   string `query:"nroHistoria"`
	NroCuenta            string `query:"nroCuenta"`
	Paciente             string `query:"paciente"`
	IdServicio           string `query:"idServicio"`
	Estado               string `query:"estado"`
	Triaje               string `query:"triaje"`
	FuenteFinanciamiento string `query:"fuenteFinanciamiento"`
}

// EstadisticasAtencion representa las métricas del dashboard
type EstadisticasAtencion struct {
	TotalCitas          int `json:"totalCitas"`
	CitasPendientes     int `json:"citasPendientes"`
	CitasAtendidas      int `json:"citasAtendidas"`
	TriajesPendientes   int `json:"triajesPendientes"`
	TriajesCompletados  int `json:"triajesCompletados"`
}
