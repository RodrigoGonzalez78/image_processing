package main

import (
	"log"
	"net/http"
	"os"

	"github.com/RodrigoGonzalez78/db"
	"github.com/RodrigoGonzalez78/middlewares"
	"github.com/RodrigoGonzalez78/routes"
	"github.com/gorilla/mux"
)

func main() {

	db.StartDB()
	db.MigrateModels()
	r := mux.NewRouter()

	r.HandleFunc("/images/{rest:.*}", routes.ServeImage).Methods("GET")

	r.HandleFunc("/register", routes.Register).Methods("POST")
	r.HandleFunc("/login", routes.Login).Methods("POST")
	r.HandleFunc("/upload", middlewares.CheckJwt(routes.Upload)).Methods("POST")

	r.HandleFunc("/user-images", middlewares.CheckJwt(routes.AllUserImages)).Methods("GET")
	r.HandleFunc("/image/{id}", middlewares.CheckJwt(routes.GetImage)).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Servidor iniciado en el puerto: " + port)
	http.ListenAndServe(":"+port, r)
}
