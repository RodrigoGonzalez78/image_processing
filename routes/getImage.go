package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/RodrigoGonzalez78/db"
	"github.com/RodrigoGonzalez78/models"
	"github.com/gorilla/mux"
)

func GetImage(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	imageIDStr := vars["id"]
	imageID, err := strconv.ParseInt(imageIDStr, 10, 64)
	if err != nil {
		http.Error(w, "ID de imagen inválido", http.StatusBadRequest)
		return
	}

	userData, ok := r.Context().Value("userData").(*models.Claim)
	if !ok || userData == nil {
		http.Error(w, "Error al obtener los datos del usuario", http.StatusUnauthorized)
	}

	image, err := db.GetImageByID(imageID)
	if err != nil {
		http.Error(w, "Imagen no encontrada: "+err.Error(), http.StatusNotFound)
		return
	}

	if image.UserName != userData.UserName {
		http.Error(w, "No tienes permiso para ver esta imagen", http.StatusUnauthorized)
		return
	}

	imageURL := fmt.Sprintf("http://localhost:8080/images/%s/%s", userData.UserName, image.Name)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"url":    imageURL,
		"name":   image.Name,
		"size":   image.Size,
		"format": image.Format,
		"width":  image.Width,
		"height": image.Height,
	})
}
