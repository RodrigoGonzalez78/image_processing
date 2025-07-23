package routes

import (
	"encoding/json"
	"image"
	"image/color"
	"net/http"
	"strconv"

	"github.com/RodrigoGonzalez78/db"
	"github.com/RodrigoGonzalez78/models"
	"github.com/RodrigoGonzalez78/storage"
	"github.com/disintegration/imaging"
	"github.com/gorilla/mux"
)

// TransformImage godoc
// @Summary      Aplica transformaciones a una imagen
// @Description  Aplica transformaciones (resize, crop, rotación, filtros) a una imagen previamente cargada por el usuario
// @Tags         images
// @Security     BearerAuth
// @Produce      image/png
// @Produce      image/jpeg
// @Produce      image/gif
// @Param        id path int true "ID de la imagen"
// @Param        body body models.TransformationRequest true "Parámetros de transformación"
// @Success      200 {file} file "Imagen transformada"
// @Failure      400 {object} models.ErrorResponse
// @Failure      404 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /images/{id}/transform [post]
func TransformImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageID := vars["id"]

	id, err := strconv.ParseInt(imageID, 10, 64)
	if err != nil {
		http.Error(w, "ID de imagen inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Obtener metadata desde la base de datos
	imageData, err := db.GetImageByID(id)
	if err != nil {
		http.Error(w, "Imagen no encontrada: "+err.Error(), http.StatusNotFound)
		return
	}

	// Obtener la imagen desde MinIO
	objReader, err := storage.MinioClientInstance.GetImage(r.Context(), imageData.Path)
	if err != nil {
		http.Error(w, "No se pudo leer la imagen desde MinIO: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer objReader.Close()

	// Decodificar imagen desde io.Reader
	srcImage, _, err := image.Decode(objReader)
	if err != nil {
		http.Error(w, "Error al decodificar la imagen: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Decodificar el JSON del body
	var transformReq models.TransformationRequest
	if err := json.NewDecoder(r.Body).Decode(&transformReq); err != nil {
		http.Error(w, "Error al parsear JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	dstImage := srcImage

	// Aplicar transformaciones
	if transformReq.Transformations.Resize.Width > 0 && transformReq.Transformations.Resize.Height > 0 {
		dstImage = imaging.Resize(dstImage, transformReq.Transformations.Resize.Width, transformReq.Transformations.Resize.Height, imaging.Lanczos)
	}

	if transformReq.Transformations.Crop.Width > 0 && transformReq.Transformations.Crop.Height > 0 {
		cropRect := image.Rect(
			transformReq.Transformations.Crop.X,
			transformReq.Transformations.Crop.Y,
			transformReq.Transformations.Crop.X+transformReq.Transformations.Crop.Width,
			transformReq.Transformations.Crop.Y+transformReq.Transformations.Crop.Height,
		)
		dstImage = imaging.Crop(dstImage, cropRect)
	}

	if transformReq.Transformations.Rotate != 0 {
		dstImage = imaging.Rotate(dstImage, transformReq.Transformations.Rotate, color.Transparent)
	}

	if transformReq.Transformations.Filters.Grayscale {
		dstImage = imaging.Grayscale(dstImage)
	}

	if transformReq.Transformations.Filters.Sepia {
		dstImage = imaging.AdjustSaturation(dstImage, -100)
		dstImage = imaging.AdjustContrast(dstImage, 10)
	}

	// Elegir formato de salida
	format := transformReq.Transformations.Format
	var imgFormat imaging.Format
	switch format {
	case "jpg", "jpeg":
		imgFormat = imaging.JPEG
		w.Header().Set("Content-Type", "image/jpeg")
	case "gif":
		imgFormat = imaging.GIF
		w.Header().Set("Content-Type", "image/gif")
	default:
		imgFormat = imaging.PNG
		w.Header().Set("Content-Type", "image/png")
	}

	// Codificar imagen en la respuesta
	w.WriteHeader(http.StatusOK)
	if err := imaging.Encode(w, dstImage, imgFormat); err != nil {
		http.Error(w, "Error al codificar la imagen: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
