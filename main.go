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

	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("uploads/"))))

	r.HandleFunc("/register", routes.Register)
	r.HandleFunc("/login", routes.Login)
	r.HandleFunc("/upload", middlewares.CheckJwt(routes.Upload))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Servidor iniciado en el puerto: " + port)
	http.ListenAndServe(":"+port, r)
}
