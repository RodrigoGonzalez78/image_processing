package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	JWTSecret      string
	BaseURL        string
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MinioUseSSL    bool
}

var Cnf *Config

func getEnv(key, fallback string, required bool) string {
	value := os.Getenv(key)
	if value == "" {
		if required {
			log.Fatalf("❌ La variable de entorno %s es obligatoria pero no está definida", key)
		}
		if fallback != "" {
			log.Printf("⚠️  La variable de entorno %s no está definida. Usando valor por defecto: %s", key, fallback)
		}
		return fallback
	}
	return value
}

func getEnvBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("⚠️  No se pudo parsear %s como boolean, usando valor por defecto: %t", key, fallback)
		return fallback
	}

	return boolValue
}

func LoadConfig() {
	// Intentar cargar .env solo en desarrollo
	if err := godotenv.Load(); err != nil {
		log.Println("📄 Archivo .env no encontrado. Usando variables de entorno del sistema.")
	}

	Cnf = &Config{
		Port:           getEnv("PORT", "8080", false),
		JWTSecret:      getEnv("JWT_SECRET", "", true), // JWT_SECRET debería ser obligatorio
		BaseURL:        getEnv("BASE_URL", "http://localhost:8080", false),
		MinioEndpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000", false),
		MinioAccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin", false),
		MinioSecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin", false),
		MinioBucket:    getEnv("MINIO_BUCKET", "images", false),
		MinioUseSSL:    getEnvBool("MINIO_USE_SSL", false),
	}

	// Validación adicional
	if Cnf.JWTSecret == "jwt_secret_key" {
		log.Println("⚠️  ADVERTENCIA: Estás usando el JWT_SECRET por defecto. Cambia esto en producción.")
	}

	log.Printf("✅ Configuración cargada exitosamente")
	log.Printf("🔧 Puerto: %s", Cnf.Port)
	log.Printf("🔧 Minio Endpoint: %s", Cnf.MinioEndpoint)
	log.Printf("🔧 Minio Bucket: %s", Cnf.MinioBucket)
	log.Printf("🔧 Minio SSL: %t", Cnf.MinioUseSSL)
}
