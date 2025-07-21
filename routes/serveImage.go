package routes

import (
	"net/http"
	"os"
	"path/filepath"
)

// ServeImage godoc
// @Summary      Servir imagen
// @Description  Devuelve el archivo de imagen almacenado en el servidor desde la carpeta local "uploads", sin exponer su ruta real.
// @Tags         images
// @Produce      image/jpeg
// @Produce      image/png
// @Produce      image/gif
// @Param        rest path string true "Ruta relativa de la imagen (por ejemplo: user123/imagen.jpg)"
// @Success      200 {file} file
// @Failure      403 {string} string "Acceso prohibido"
// @Failure      404 {string} string "Archivo no encontrado"
// @Router       /images/{rest} [get]
func ServeImage(w http.ResponseWriter, r *http.Request) {

	filePath := filepath.Join("uploads", r.URL.Path[len("/images/"):])

	info, err := os.Stat(filePath)
	if err != nil {
		http.Error(w, "Archivo o directorio no encontrado", http.StatusNotFound)
		return
	}
	if info.IsDir() {
		http.Error(w, "Acceso prohibido", http.StatusForbidden)
		return
	}

	http.ServeFile(w, r, filePath)
}
