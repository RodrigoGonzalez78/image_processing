package main

import (
	"github.com/RodrigoGonzalez78/routes"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/upload", routes.Upload)
}
