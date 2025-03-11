package routes

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/RodrigoGonzalez78/db"
	"github.com/RodrigoGonzalez78/models"
	"github.com/google/uuid"
)

func Upload(w http.ResponseWriter, r *http.Request) {

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error al obtener el archivo: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	userData, ok := r.Context().Value("userData").(*models.Claim)
	if !ok || userData == nil {
		http.Error(w, "Error al obtener los datos del usuario", http.StatusUnauthorized)
		return
	}

	extParts := strings.Split(handler.Filename, ".")
	if len(extParts) < 2 {
		http.Error(w, "Archivo sin extensión válida", http.StatusBadRequest)
		return
	}
	extension := extParts[len(extParts)-1]

	directory := "uploads/" + userData.UserName
	err = os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		http.Error(w, "Error al crear el directorio: "+err.Error(), http.StatusInternalServerError)
		return
	}

	randomName := uuid.New().String() + "." + extension

	filename := filepath.Join(directory, randomName)

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, "Error al abrir el archivo: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, "Error al copiar la imagen: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.CreateImage(randomName, userData.UserName)
	if err != nil {
		http.Error(w, "Error al almacenar la descripción de la imagen en la BD: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Imagen subida exitosamente",
		"image":   randomName,
	})
}
