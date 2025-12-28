package model

import "sihce_consulta_externa/internal/config/database"

// AppConfig contiene toda la configuración de la aplicación
type AppConfig struct {
	JWT      JWTConfig      `toml:"jwt"`
	Security SecurityConfig `toml:"security"`
	App      ServerConfig   `toml:"app"`
	Database database.DatabaseConfig `toml:"database"`
}

// JWTConfig contiene la configuración de JWT
type JWTConfig struct {
	AccessSecret                 string `toml:"access_secret"`
	RefreshSecret                string `toml:"refresh_secret"`
	AccessTokenExpirationSeconds int    `toml:"access_token_expiration_seconds"`
	RefreshTokenExpirationSeconds int   `toml:"refresh_token_expiration_seconds"`
}

// SecurityConfig contiene la configuración de seguridad
type SecurityConfig struct {
	SessionSecret  string `toml:"session_secret"`
	HashSaltRounds int    `toml:"hash_salt_rounds"`
}

// ServerConfig contiene la configuración del servidor
type ServerConfig struct {
	Port   int    `toml:"port"`
	AppEnv string `toml:"app_env"`
}
