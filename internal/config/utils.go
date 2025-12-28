package config

import (
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

// ExpandEnvVars expande las variables de entorno en el formato ${VAR_NAME}
func ExpandEnvVars(input string) string {
	re := regexp.MustCompile(`\$\{([^}]+)\}`)
	return re.ReplaceAllStringFunc(input, func(match string) string {
		varName := match[2 : len(match)-1] // Extrae el nombre de la variable
		if value := os.Getenv(varName); value != "" {
			return value
		}
		return match // Si no existe, mantiene el valor original
	})
}

// CargarEnv carga el archivo .env apropiado según el entorno del config.toml
func CargarEnv(appEnv string) error {
	if appEnv == "" {
		appEnv = "development" // Por defecto usar development
	}

	var envFile string

	// Determinar el archivo .env a cargar según el entorno
	switch appEnv {
	case "development", "dev":
		envFile = ".env.dev"
	case "production", "prod":
		envFile = ".env.prod"
	default:
		envFile = ".env.dev"
	}

	// Intentar cargar el archivo .env específico del entorno
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Advertencia: No se encontró archivo %s, usando variables de entorno del sistema", envFile)
		return nil
	}

	return nil
}