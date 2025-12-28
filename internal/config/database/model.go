package database

type PoolConfig struct {
	Min                 int `toml:"min"`
	Max                 int `toml:"max"`
	IdleTimeoutMs       int `toml:"idle_timeout_ms"`
	ConnectionTimeoutMs int `toml:"connection_timeout_ms"`
}

type DBConfig struct {
	Host            string          `toml:"host"`
	Port            int             `toml:"port"`
	Name            string          `toml:"name"`
	User            string          `toml:"user"`
	Password        string          `toml:"password"`
	Encrypt         bool            `toml:"encrypt"`
	TrustServerCert bool            `toml:"trust_server_certificate"`
	Pool            PoolConfig `toml:"pool"`
}

type ConfigDatabase struct {
	Database DatabaseConfig `toml:"database"`
}

type DatabaseConfig struct {
	Principal  DBConfig `toml:"principal"`
	Secundaria DBConfig `toml:"secundaria"`
}