package routes

import (
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/RodrigoGonzalez78/storage"
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

	objectPath := strings.TrimPrefix(r.URL.Path, "/images/")
	objectPath = path.Clean(objectPath)

	if objectPath == "." || strings.Contains(objectPath, "..") {
		http.Error(w, "Acceso prohibido", http.StatusForbidden)
		return
	}

	// Obtener objeto desde MinIO
	reader, err := storage.MinioClientInstance.GetImage(r.Context(), objectPath)
	if err != nil {
		http.Error(w, "Archivo no encontrado en MinIO", http.StatusNotFound)
		return
	}
	defer reader.Close()

	// Estimar tipo MIME en base a la extensión
	switch ext := strings.ToLower(path.Ext(objectPath)); ext {
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".gif":
		w.Header().Set("Content-Type", "image/gif")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	w.WriteHeader(http.StatusOK)
	io.Copy(w, reader)
}
