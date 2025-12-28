package atenciones

import (
	"context"
	"database/sql"
	"time"

	"sihce_consulta_externa/internal/config/database"
	http_atenciones "sihce_consulta_externa/internal/modules/atenciones/http"
)

// Repository interfaz para el repositorio de atenciones
type Repository interface {
	ListarAtenciones(ctx context.Context, filtros *http_atenciones.FiltrosAtencion) ([]http_atenciones.AtencionRepositoryDTO, error)
	ObtenerAtencionPorId(ctx context.Context, idAtencion int) (*http_atenciones.AtencionRepositoryDTO, error)
	ObtenerEstadisticas(ctx context.Context, fechaInicio, fechaFin string) (*http_atenciones.EstadisticasAtencion, error)
}

// RepositoryImpl implementación del repositorio
type RepositoryImpl struct {
	db *database.GestorDB
}

// NuevoRepository crea una nueva instancia del repositorio
func NuevoRepository(db *database.GestorDB) Repository {
	return &RepositoryImpl{db: db}
}

// ListarAtenciones obtiene todas las atenciones con filtros opcionales
func (r *RepositoryImpl) ListarAtenciones(ctx context.Context, filtros *http_atenciones.FiltrosAtencion) ([]http_atenciones.AtencionRepositoryDTO, error) {
	conn, err := r.db.GetPrincipal()
	if err != nil {
		return nil, err
	}

	// Por ahora usamos el query mock, pero aquí irá la query real con filtros
	query := QUERY_ATENCIONES

	rows, err := conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var atencionesList []http_atenciones.AtencionRepositoryDTO

	for rows.Next() {
		var a http_atenciones.AtencionRepositoryDTO
		var fechaCitaStr string

		err := rows.Scan(
			&a.IdPaciente,
			&a.NroDocumento,
			&a.FuenteFinanciamiento,
			&a.IdAtencion,
			&a.NombreDoctor,
			&a.ApellidoDoctor,
			&a.ApellidosPaciente,
			&a.NombresPaciente,
			&a.Edad,
			&a.NroHistoriaClinica,
			&a.Servicio,
			&a.IdServicio,
			&fechaCitaStr,
			&a.HoraInicio,
			&a.HoraFin,
			&a.Estado,
			&a.Triaje,
		)
		if err != nil {
			return nil, err
		}

		// Parsear la fecha del string
		fechaCita, err := time.Parse("2006-01-02", fechaCitaStr)
		if err != nil {
			// Si falla el parseo, usar fecha actual
			fechaCita = time.Now()
		}
		a.FechaCita = fechaCita

		atencionesList = append(atencionesList, a)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return atencionesList, nil
}

// ObtenerAtencionPorId obtiene una atención por su ID
func (r *RepositoryImpl) ObtenerAtencionPorId(ctx context.Context, idAtencion int) (*http_atenciones.AtencionRepositoryDTO, error) {
	conn, err := r.db.GetPrincipal()
	if err != nil {
		return nil, err
	}

	// Query para obtener una atención específica
	query := `
		SELECT
			IdPaciente, NroDocumento, FuenteFinanciamiento, IdAtencion,
			NombreDoctor, ApellidoDoctor, ApellidosPaciente, NombresPaciente,
			Edad, NroHistoriaClinica, Servicio, IdServicio,
			FechaCita, HoraInicio, HoraFin, Estado, Triaje
		FROM (` + QUERY_ATENCIONES + `) AS atenciones
		WHERE IdAtencion = @p1
	`

	row := conn.QueryRowContext(ctx, query, idAtencion)

	var a http_atenciones.AtencionRepositoryDTO
	var fechaCitaStr string

	err = row.Scan(
		&a.IdPaciente,
		&a.NroDocumento,
		&a.FuenteFinanciamiento,
		&a.IdAtencion,
		&a.NombreDoctor,
		&a.ApellidoDoctor,
		&a.ApellidosPaciente,
		&a.NombresPaciente,
		&a.Edad,
		&a.NroHistoriaClinica,
		&a.Servicio,
		&a.IdServicio,
		&fechaCitaStr,
		&a.HoraInicio,
		&a.HoraFin,
		&a.Estado,
		&a.Triaje,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No encontrado
		}
		return nil, err
	}

	// Parsear la fecha del string
	fechaCita, err := time.Parse("2006-01-02", fechaCitaStr)
	if err != nil {
		fechaCita = time.Now()
	}
	a.FechaCita = fechaCita

	return &a, nil
}

// ObtenerEstadisticas calcula las estadísticas del dashboard
func (r *RepositoryImpl) ObtenerEstadisticas(ctx context.Context, fechaInicio, fechaFin string) (*http_atenciones.EstadisticasAtencion, error) {
	conn, err := r.db.GetPrincipal()
	if err != nil {
		return nil, err
	}

	// Query simplificada: calcular estadísticas directamente desde la query base
	query := `
		WITH atenciones AS (
			` + QUERY_ATENCIONES + `
		)
		SELECT
			COUNT(*) as TotalCitas,
			ISNULL(SUM(CASE WHEN Estado = 'Pendiente' THEN 1 ELSE 0 END), 0) as CitasPendientes,
			ISNULL(SUM(CASE WHEN Estado = 'Atendido' THEN 1 ELSE 0 END), 0) as CitasAtendidas,
			ISNULL(SUM(CASE WHEN Triaje IN ('Pendiente', 'No aplica') THEN 1 ELSE 0 END), 0) as TriajesPendientes,
			ISNULL(SUM(CASE WHEN Triaje = 'Completado' THEN 1 ELSE 0 END), 0) as TriajesCompletados
		FROM atenciones
	`

	var stats http_atenciones.EstadisticasAtencion

	err = conn.QueryRowContext(ctx, query).Scan(
		&stats.TotalCitas,
		&stats.CitasPendientes,
		&stats.CitasAtendidas,
		&stats.TriajesPendientes,
		&stats.TriajesCompletados,
	)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
