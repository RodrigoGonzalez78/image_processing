package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/RodrigoGonzalez78/db"
	"github.com/RodrigoGonzalez78/models"
)

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

	// Construir la respuesta con la información de paginación
	response := map[string]interface{}{
		"page":   page,
		"limit":  limit,
		"total":  total,
		"images": images,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
