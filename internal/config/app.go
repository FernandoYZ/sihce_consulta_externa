package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"

	"sihce_consulta_externa/internal/config/database"
	"sihce_consulta_externa/internal/model"
)

var AppConfig *model.AppConfig

// CargarConfiguracion carga la configuración desde archivos TOML y variables de entorno
func CargarConfiguracion() error {
	// 1. Leer config.toml para obtener app_env
	appConfigPath := filepath.Join("internal", "config", "config.toml")
	appConfigBytes, err := os.ReadFile(appConfigPath)
	if err != nil {
		return fmt.Errorf("error leyendo config.toml: %w", err)
	}

	// Parsear solo app_env (sin expansión de variables)
	var tempConfig struct {
		App struct {
			AppEnv string `toml:"app_env"`
		} `toml:"app"`
	}
	if err := toml.Unmarshal(appConfigBytes, &tempConfig); err != nil {
		return fmt.Errorf("error parseando app_env: %w", err)
	}

	// 2. Cargar archivo .env correspondiente
	if err := CargarEnv(tempConfig.App.AppEnv); err != nil {
		return fmt.Errorf("error cargando variables de entorno: %w", err)
	}

	// 3. Leer database.toml
	dbConfigPath := filepath.Join("internal", "config", "database", "database.toml")
	dbConfigBytes, err := os.ReadFile(dbConfigPath)
	if err != nil {
		return fmt.Errorf("error leyendo database.toml: %w", err)
	}

	// 4. Expandir variables de entorno en ambos archivos
	appConfigContent := ExpandEnvVars(string(appConfigBytes))
	dbConfigContent := ExpandEnvVars(string(dbConfigBytes))

	// 5. Parsear configuración principal
	AppConfig = &model.AppConfig{}
	if err := toml.Unmarshal([]byte(appConfigContent), AppConfig); err != nil {
		return fmt.Errorf("error parseando config.toml: %w", err)
	}

	// 6. Parsear y asignar configuración de base de datos
	var dbConfig struct {
		Database database.DatabaseConfig `toml:"database"`
	}
	if err := toml.Unmarshal([]byte(dbConfigContent), &dbConfig); err != nil {
		return fmt.Errorf("error parseando database.toml: %w", err)
	}
	AppConfig.Database = dbConfig.Database

	return nil
}

// ObtenerConfiguracion devuelve la configuración de la aplicación
func ObtenerConfiguracion() *model.AppConfig {
	if AppConfig == nil {
		log.Fatal("La configuración no ha sido cargada. Llama a CargarConfiguracion() primero")
	}
	return AppConfig
}
