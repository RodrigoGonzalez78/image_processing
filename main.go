package main

import (
	"github.com/RodrigoGonzalez78/db"
	"github.com/RodrigoGonzalez78/routes"
	"github.com/gorilla/mux"
)

func main() {

	db.StartDB()
	db.MigrateModels()
	r := mux.NewRouter()

	r.HandleFunc("/register", routes.Register)
	r.HandleFunc("/login", routes.Login)
	r.HandleFunc("/upload", routes.Upload)
}
