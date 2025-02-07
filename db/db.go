package db

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func StartDB() {
	db, err := gorm.Open(sqlite.Open("mi_gorm_sqlite.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Fallo al conectar a la base de datos SQLite:", err)
	}
	fmt.Println("Conexión a la base de datos SQLite (GORM) establecida.")

	db.AutoMigrate(&Usuario{})
	fmt.Println("Esquema de la base de datos migrado (tabla 'users' creada/actualizada).")
}
