package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"

	// Config
	"github.com/RodrigoGonzalez78/config"

	// Domain services
	"github.com/RodrigoGonzalez78/internal/domain/service"

	// Application use cases
	authUC "github.com/RodrigoGonzalez78/internal/application/usecase/auth"
	imageUC "github.com/RodrigoGonzalez78/internal/application/usecase/image"

	// Infrastructure
	gormDB "github.com/RodrigoGonzalez78/internal/infrastructure/persistence/gorm"
	minioStorage "github.com/RodrigoGonzalez78/internal/infrastructure/storage/minio"

	// HTTP
	"github.com/RodrigoGonzalez78/internal/infrastructure/http/handler"
	"github.com/RodrigoGonzalez78/internal/infrastructure/http/middleware"

	// Swagger
	"github.com/RodrigoGonzalez78/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           API de Procesamiento de Imágenes
// @version         1.0
// @description     Esta es una API para subir, transformar y consultar imágenes.
// @termsOfService  http://swagger.io/terms/

// @BasePath        /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Escribe "Bearer " seguido del token JWT obtenido en /login
func main() {

	config.LoadConfig()

	baseURL := strings.TrimPrefix(config.Cnf.BaseURL, "http://")
	baseURL = strings.TrimPrefix(baseURL, "https://")
	docs.SwaggerInfo.Host = baseURL + ":" + config.Cnf.Port

	database, err := gormDB.NewDatabase("database.db")
	if err != nil {
		log.Fatal(" Error al conectar a la base de datos:", err)
	}

	if err := database.Migrate(); err != nil {
		log.Fatal(" Error al migrar la base de datos:", err)
	}

	fileStorage, err := minioStorage.NewFileStorage(
		config.Cnf.MinioEndpoint,
		config.Cnf.MinioAccessKey,
		config.Cnf.MinioSecretKey,
		config.Cnf.MinioBucket,
		config.Cnf.MinioUseSSL,
	)
	if err != nil {
		log.Fatal(" Error al conectar a MinIO:", err)
	}

	userRepo := gormDB.NewUserRepository(database.DB)
	imageRepo := gormDB.NewImageRepository(database.DB)

	passwordService := service.NewPasswordService()
	tokenService := service.NewTokenService(config.Cnf.JWTSecret, 24*time.Hour)

	registerUC := authUC.NewRegisterUseCase(userRepo, passwordService)
	loginUC := authUC.NewLoginUseCase(userRepo, passwordService, tokenService)

	uploadUC := imageUC.NewUploadUseCase(imageRepo, fileStorage, config.Cnf.BaseURL, config.Cnf.Port)
	getImageUC := imageUC.NewGetUseCase(imageRepo, config.Cnf.BaseURL, config.Cnf.Port)
	listImagesUC := imageUC.NewListUseCase(imageRepo, config.Cnf.BaseURL, config.Cnf.Port)
	transformUC := imageUC.NewTransformUseCase(imageRepo, fileStorage)

	authHandler := handler.NewAuthHandler(registerUC, loginUC)
	imageHandler := handler.NewImageHandler(uploadUC, getImageUC, listImagesUC, transformUC)

	jwtMiddleware := middleware.NewJWTMiddleware(tokenService)

	r := mux.NewRouter()

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	r.HandleFunc("/images/{rest:.*}", imageHandler.ServeImage).Methods("GET")

	r.HandleFunc("/upload", jwtMiddleware.Authenticate(imageHandler.Upload)).Methods("POST")
	r.HandleFunc("/user-images", jwtMiddleware.Authenticate(imageHandler.ListUserImages)).Methods("GET")
	r.HandleFunc("/images/{id}", jwtMiddleware.Authenticate(imageHandler.GetImage)).Methods("GET")
	r.HandleFunc("/images/{id}/transform", jwtMiddleware.Authenticate(imageHandler.TransformImage)).Methods("POST")

	log.Println(" Servidor iniciado en el puerto:", config.Cnf.Port)
	log.Println(" Swagger UI disponible en: http://localhost:" + config.Cnf.Port + "/swagger/index.html")
	http.ListenAndServe(":"+config.Cnf.Port, r)
}
