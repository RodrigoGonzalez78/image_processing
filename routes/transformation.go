package routes

import (
	"encoding/json"
	"image"
	"image/color"
	"net/http"
	"strconv"

	"github.com/RodrigoGonzalez78/db"
	"github.com/RodrigoGonzalez78/models"
	"github.com/disintegration/imaging"
	"github.com/gorilla/mux"
)

func TransformImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageID := vars["id"]

	id, err := strconv.ParseInt(imageID, 10, 64)
	if err != nil {
		http.Error(w, "ID de imagen inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	imageData, err := db.GetImageByID(id)
	if err != nil {
		http.Error(w, "Imagen no encontrada: "+err.Error(), http.StatusNotFound)
		return
	}

	srcImage, err := imaging.Open(imageData.Path)
	if err != nil {
		http.Error(w, "Error al abrir la imagen: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var transformReq models.TransformationRequest
	if err := json.NewDecoder(r.Body).Decode(&transformReq); err != nil {
		http.Error(w, "Error al parsear JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	dstImage := srcImage

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
		// Se simula el efecto sepia con ajustes de saturación, contraste y hue.
		dstImage = imaging.AdjustSaturation(dstImage, -100)
		dstImage = imaging.AdjustContrast(dstImage, 10)
		// Usar otra librería para el ajuste de matiz
		// o eliminar esta transformación.
		// dstImage = imaging.AdjustHue(dstImage, 30)
	}

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

	// Codificar y escribir la imagen transformada directamente a la respuesta
	if err := imaging.Encode(w, dstImage, imgFormat); err != nil {
		http.Error(w, "Error al codificar la imagen: "+err.Error(), http.StatusInternalServerError)
		return
	}

}
