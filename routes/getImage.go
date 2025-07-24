package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/RodrigoGonzalez78/config"
	"github.com/RodrigoGonzalez78/db"
	"github.com/RodrigoGonzalez78/models"
	"github.com/gorilla/mux"
)

// GetImage godoc
// @Summary      Obtener información de una imagen
// @Description  Devuelve la metadata y URL de una imagen del usuario autenticado
// @Tags         images
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "ID de la imagen"
// @Success      200 {object} models.ImageDetailResponse
// @Failure      400 {object} models.ErrorResponse
// @Failure      401 {object} models.ErrorResponse
// @Failure      404 {object} models.ErrorResponse
// @Router       /image/{id} [get]
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
		http.Error(w, "No autorizado", http.StatusUnauthorized)
		return
	}

	image, err := db.GetImageByID(imageID)
	if err != nil {
		http.Error(w, "Imagen no encontrada", http.StatusNotFound)
		return
	}

	if image.UserName != userData.UserName {
		http.Error(w, "No tienes permiso para esta imagen", http.StatusUnauthorized)
		return
	}

	imageURL := fmt.Sprintf("%s:%s/%s/%s/%s", config.Cnf.BaseURL, config.Cnf.Port, "images", image.UserName, image.Name)

	resp := models.ImageDetailResponse{
		URL:    imageURL,
		Name:   image.Name,
		Size:   image.Size,
		Format: image.Format,
		Width:  image.Width,
		Height: image.Height,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
