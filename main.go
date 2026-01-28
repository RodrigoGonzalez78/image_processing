package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/RodrigoGonzalez78/config"
	"github.com/RodrigoGonzalez78/db"
	"github.com/RodrigoGonzalez78/docs"
	"github.com/RodrigoGonzalez78/middlewares"
	"github.com/RodrigoGonzalez78/routes"
	"github.com/RodrigoGonzalez78/storage"
	"github.com/gorilla/mux"

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

	db.StartDB()
	db.MigrateModels()
	config.LoadConfig()

	// Configurar Swagger dinámicamente con la URL y puerto actual
	baseURL := strings.TrimPrefix(config.Cnf.BaseURL, "http://")
	baseURL = strings.TrimPrefix(baseURL, "https://")
	docs.SwaggerInfo.Host = baseURL + ":" + config.Cnf.Port

	err := storage.StartMinioClient(
		config.Cnf.MinioEndpoint,
		config.Cnf.MinioAccessKey,
		config.Cnf.MinioSecretKey,
		config.Cnf.MinioBucket,
		config.Cnf.MinioUseSSL,
	)

	if err != nil {
		log.Fatal("Fallo al conectar a Minio:", err)

	}

	r := mux.NewRouter()

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	r.HandleFunc("/images/{rest:.*}", routes.ServeImage).Methods("GET")

	r.HandleFunc("/register", routes.Register).Methods("POST")
	r.HandleFunc("/login", routes.Login).Methods("POST")
	r.HandleFunc("/upload", middlewares.CheckJwt(routes.Upload)).Methods("POST")

	r.HandleFunc("/user-images", middlewares.CheckJwt(routes.AllUserImages)).Methods("GET")
	r.HandleFunc("/images/{id}", middlewares.CheckJwt(routes.GetImage)).Methods("GET")
	r.HandleFunc("/images/{id}/transform", middlewares.CheckJwt(routes.TransformImage)).Methods("POST")

	log.Println("Servidor iniciado en el puerto: " + config.Cnf.Port)
	http.ListenAndServe(":"+config.Cnf.Port, r)
}
