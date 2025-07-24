package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/RodrigoGonzalez78/config"
	"github.com/RodrigoGonzalez78/db"
	"github.com/RodrigoGonzalez78/models"
)

// AllUserImages godoc
// @Summary      Lista imágenes del usuario autenticado
// @Description  Devuelve las imágenes subidas por el usuario autenticado con soporte de paginación.
// @Tags         images
// @Security     BearerAuth
// @Produce      json
// @Param        page   query int false "Número de página" default(1)
// @Param        limit  query int false "Cantidad por página" default(10)
// @Success      200 {object} models.PaginatedImagesResponse
// @Failure      401 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /user-images [get]
func AllUserImages(w http.ResponseWriter, r *http.Request) {

	userData, ok := r.Context().Value("userData").(*models.Claim)
	if !ok || userData == nil {
		http.Error(w, "Usuario no autenticado", http.StatusUnauthorized)
		return
	}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	images, total, err := db.GetImagesByUser(userData.UserName, page, limit)
	if err != nil {
		http.Error(w, "Error al obtener imágenes: "+err.Error(), http.StatusInternalServerError)
		return
	}
	imagesWIthURL := make([]models.Image, len(images))
	for i, img := range images {
		img.Path = fmt.Sprintf("%s:%s/%s/%s/%s", config.Cnf.BaseURL, config.Cnf.Port, "images", img.UserName, img.Name)
		imagesWIthURL[i] = img
	}

	// Construir la respuesta con la información de paginación
	response := models.PaginatedImagesResponse{
		Page:   page,
		Limit:  limit,
		Total:  total,
		Images: imagesWIthURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
