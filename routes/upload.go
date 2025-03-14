package routes

import (
	"encoding/json"
	"fmt"
	"image"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

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

	allowedExtensions := map[string]bool{"jpg": true, "jpeg": true, "png": true, "gif": true}
	if !allowedExtensions[extension] {
		http.Error(w, "Formato de archivo no permitido", http.StatusBadRequest)
		return
	}

	directory := filepath.Join("uploads", userData.UserName)
	err = os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		http.Error(w, "Error al crear el directorio: "+err.Error(), http.StatusInternalServerError)
		return
	}

	randomName := uuid.New().String() + "." + extension
	filePath := filepath.Join(directory, randomName)

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, "Error al abrir el archivo: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	size, err := io.Copy(f, file)
	if err != nil {
		http.Error(w, "Error al copiar la imagen: "+err.Error(), http.StatusInternalServerError)
		return
	}

	imageFile, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Error al abrir la imagen para analizar: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer imageFile.Close()

	img, format, err := image.DecodeConfig(imageFile)
	if err != nil {
		http.Error(w, "Error al obtener dimensiones de la imagen: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Guardar en la base de datos
	imageData := models.Image{
		Name:     randomName,
		UserName: userData.UserName,
		Path:     filePath,
		Size:     size,
		Format:   format,
		Width:    img.Width,
		Height:   img.Height,
	}

	err = db.CreateImage(imageData)

	if err != nil {
		http.Error(w, "Error al almacenar la imagen en la base de datos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	imageURL := fmt.Sprintf("http://localhost:8080/images/%s/%s", userData.UserName, randomName)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Imagen subida exitosamente",
		"image": map[string]interface{}{
			"url":    imageURL,
			"name":   randomName,
			"size":   size,
			"format": format,
			"width":  img.Width,
			"height": img.Height,
		},
	})
}
