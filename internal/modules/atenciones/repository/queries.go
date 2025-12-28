package atenciones

const (
	QUERY_ATENCIONES = `
		select 
			141931 as IdPaciente,
			'70059943' as NroDocumento,
			'ESTRATEGIA' as FuenteFinanciamiento,
			402971 as IdAtencion,
			'zaida angelica' as NombreDoctor,
			'soto reyes' as ApellidoDoctor,
			'barraza meza' as ApellidosPaciente,
			'angello gabriel' as NombresPaciente,
			28 as Edad,
			95734 as NroHistoriaClinica,
			'terapia fisica 4 - tce.' as Servicio,
			414 as IdServicio,
			'2025-12-23' as FechaCita,
			'07:00' as HoraInicio,
			'07:30' as HoraFin,
			'Pendiente' as Estado,
			'No aplica' as Triaje
		union all
		select 
			393595 as IdPaciente,
			'91117761' as NroDocumento,
			'ESTRATEGIA' as FuenteFinanciamiento,
			409027 as IdAtencion,
			'martha elizabeth' as NombreDoctor,
			'osco solorzano' as ApellidoDoctor,
			'laura huaroto' as ApellidosPaciente,
			'dyland gael' as NombresPaciente,
			6 as Edad,
			304609 as NroHistoriaClinica,
			'terapia ocupacional 1 - t.ce' as Servicio,
			408 as IdServicio,
			'2025-12-23' as FechaCita,
			'07:00' as HoraInicio,
			'07:40' as HoraFin,
			'Atendido' as Estado,
			'Completado' as Triaje
		union all
		select 
			421416 as IdPaciente,
			'93126568' as NroDocumento,
			'ESTRATEGIA' as FuenteFinanciamiento,
			413313 as IdAtencion,
			'caroll rosario' as NombreDoctor,
			'anchante gonzales del valle' as ApellidoDoctor,
			'lapa villafuerte' as ApellidosPaciente,
			'john ryan' as NombresPaciente,
			3 as Edad,
			333130 as NroHistoriaClinica,
			'terapia de lenguaje - 1 .tce' as Servicio,
			409 as IdServicio,
			'2025-12-23' as FechaCita,
			'07:00' as HoraInicio,
			'07:30' as HoraFin,
			'Atendido' as Estado,
			'Completado' as Triaje
		union all
		select 
			176769 as IdPaciente,
			'15358647' as NroDocumento,
			'ESTRATEGIA' as FuenteFinanciamiento,
			413467 as IdAtencion,
			'magnolia yesenia' as NombreDoctor,
			'avalos cuya' as ApellidoDoctor,
			'julian soto' as ApellidosPaciente,
			'gilda ivana' as NombresPaciente,
			52 as Edad,
			11792 as NroHistoriaClinica,
			'terapia fisica 3 - tce.' as Servicio,
			412 as IdServicio,
			'2025-12-23' as FechaCita,
			'07:00' as HoraInicio,
			'07:30' as HoraFin,
			'Pendiente' as Estado,
			'No aplica' as Triaje
		union all
		select 
			157193 as IdPaciente,
			'15354151' as NroDocumento,
			'ESTRATEGIA' as FuenteFinanciamiento,
			413479 as IdAtencion,
			'cindy viviana' as NombreDoctor,
			'moreno suice' as ApellidoDoctor,
			'de la cruz palomino' as ApellidosPaciente,
			'juan roberto' as NombresPaciente,
			58 as Edad,
			31906 as NroHistoriaClinica,
			'terapia fisica 2 - tce.' as Servicio,
			411 as IdServicio,
			'2025-12-23' as FechaCita,
			'07:00' as HoraInicio,
			'07:30' as HoraFin,
			'Atendido' as Estado,
			'Completado' as Triaje
	`
)