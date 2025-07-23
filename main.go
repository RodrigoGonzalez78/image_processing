package main

import (
	"log"
	"net/http"

	"github.com/RodrigoGonzalez78/config"
	"github.com/RodrigoGonzalez78/db"
	"github.com/RodrigoGonzalez78/middlewares"
	"github.com/RodrigoGonzalez78/routes"
	"github.com/RodrigoGonzalez78/storage"
	"github.com/gorilla/mux"

	_ "github.com/RodrigoGonzalez78/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           API de Procesamiento de Imágenes
// @version         1.0
// @description     Esta es una API para subir, transformar y consultar imágenes.
// @host            localhost:8080
// @BasePath        /
func main() {

	db.StartDB()
	db.MigrateModels()
	config.LoadConfig()

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
	r.HandleFunc("/image/{id}", middlewares.CheckJwt(routes.GetImage)).Methods("GET")
	r.HandleFunc("/images/{id}/transform", middlewares.CheckJwt(routes.TransformImage)).Methods("POST")

	log.Println("Servidor iniciado en el puerto: " + config.Cnf.Port)
	http.ListenAndServe(":"+config.Cnf.Port, r)
}
