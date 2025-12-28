package database

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/microsoft/go-mssqldb"
)

// GestorDB gestiona los pools de conexión a las bases de datos
type GestorDB struct {
	principal     *sql.DB
	secundaria    *sql.DB
	configuracion *DatabaseConfig
	mu            sync.RWMutex
}

var (
	instancia *GestorDB
	once      sync.Once
)

// NewGestorDB crea y retorna la instancia singleton del gestor de bases de datos
func NewGestorDB(config *DatabaseConfig) (*GestorDB, error) {
	var err error

	once.Do(func() {
		instancia = &GestorDB{
			configuracion: config,
		}
	})

	return instancia, err
}

// GetGestor devuelve la instancia singleton del gestor
func GetGestor() (*GestorDB, error) {
	if instancia == nil {
		return nil, fmt.Errorf("gestor no inicializado, llamar NewGestorDB primero")
	}
	return instancia, nil
}

// construirCadenaConexion construye la cadena de conexión para SQL Server
func construirCadenaConexion(cfg DBConfig) string {
	encrypt := "disable"
	if cfg.Encrypt {
		encrypt = "true"
	}

	trustCert := "false"
	if cfg.TrustServerCert {
		trustCert = "true"
	}

	return fmt.Sprintf(
		"server=%s;port=%d;database=%s;user id=%s;password=%s;encrypt=%s;TrustServerCertificate=%s;connection timeout=%d",
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.User,
		cfg.Password,
		encrypt,
		trustCert,
		cfg.Pool.ConnectionTimeoutMs/1000,
	)
}

// inicializarPool inicializa un pool de conexiones
func inicializarPool(cfg DBConfig, nombre string) (*sql.DB, error) {
	cadenaConexion := construirCadenaConexion(cfg)

	db, err := sql.Open("sqlserver", cadenaConexion)
	if err != nil {
		return nil, fmt.Errorf("error al abrir conexión %s: %w", nombre, err)
	}

	// Configurar el pool de conexiones
	db.SetMaxOpenConns(cfg.Pool.Max)
	db.SetMaxIdleConns(cfg.Pool.Min)
	db.SetConnMaxIdleTime(time.Duration(cfg.Pool.IdleTimeoutMs) * time.Millisecond)
	db.SetConnMaxLifetime(30 * time.Minute)

	// Verificar la conexión con contexto y timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Pool.ConnectionTimeoutMs)*time.Millisecond)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("error al hacer ping a la base de datos %s: %w", nombre, err)
	}

	return db, nil
}

// GetPrincipal obtiene el pool de conexión principal (lazy initialization)
func (g *GestorDB) GetPrincipal() (*sql.DB, error) {
	// Lectura rápida sin bloqueo si ya existe
	g.mu.RLock()
	if g.principal != nil {
		g.mu.RUnlock()
		return g.principal, nil
	}
	g.mu.RUnlock()

	// Bloqueo de escritura para inicializar
	g.mu.Lock()
	defer g.mu.Unlock()

	// Double-check locking pattern
	if g.principal != nil {
		return g.principal, nil
	}

	db, err := inicializarPool(g.configuracion.Principal, "principal")
	if err != nil {
		return nil, err
	}

	g.principal = db
	return g.principal, nil
}

// GetSecundaria obtiene el pool de conexión secundaria (lazy initialization)
func (g *GestorDB) GetSecundaria() (*sql.DB, error) {
	// Lectura rápida sin bloqueo si ya existe
	g.mu.RLock()
	if g.secundaria != nil {
		g.mu.RUnlock()
		return g.secundaria, nil
	}
	g.mu.RUnlock()

	// Bloqueo de escritura para inicializar
	g.mu.Lock()
	defer g.mu.Unlock()

	// Double-check locking pattern
	if g.secundaria != nil {
		return g.secundaria, nil
	}

	db, err := inicializarPool(g.configuracion.Secundaria, "secundaria")
	if err != nil {
		return nil, err
	}

	g.secundaria = db
	return g.secundaria, nil
}

// VerificarSalud verifica el estado de las conexiones activas
func (g *GestorDB) VerificarSalud(ctx context.Context) error {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if g.principal != nil {
		if err := g.principal.PingContext(ctx); err != nil {
			return fmt.Errorf("error en health check de base de datos principal: %w", err)
		}
	}

	if g.secundaria != nil {
		if err := g.secundaria.PingContext(ctx); err != nil {
			return fmt.Errorf("error en health check de base de datos secundaria: %w", err)
		}
	}

	return nil
}

// Close cierra todas las conexiones de base de datos
func (g *GestorDB) ApagarDatabase() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	var errores []error

	if g.principal != nil {
		if err := g.principal.Close(); err != nil {
			errores = append(errores, fmt.Errorf("error al cerrar pool principal: %w", err))
		}
		g.principal = nil
	}

	if g.secundaria != nil {
		if err := g.secundaria.Close(); err != nil {
			errores = append(errores, fmt.Errorf("error al cerrar pool secundaria: %w", err))
		}
		g.secundaria = nil
	}

	if len(errores) > 0 {
		return fmt.Errorf("errores al cerrar conexiones: %v", errores)
	}

	return nil
}
