package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"net/http"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/RodrigoGonzalez78/config"
	"github.com/RodrigoGonzalez78/db"
	"github.com/RodrigoGonzalez78/models"
	"github.com/RodrigoGonzalez78/storage"
	"github.com/google/uuid"
)

// Upload godoc
// @Summary      Sube una imagen
// @Description  Permite a un usuario autenticado subir una imagen. Se guarda en MinIO y se almacena la metadata.
// @Tags         images
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        image formData file true "Imagen a subir (jpg, jpeg, png, gif)"
// @Success      201 {object} models.UploadResponse
// @Failure      400 {object} models.ErrorResponse
// @Failure      401 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /upload [post]
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

	extension := strings.ToLower(extParts[len(extParts)-1])
	allowedExtensions := map[string]bool{"jpg": true, "jpeg": true, "png": true, "gif": true}

	if !allowedExtensions[extension] {
		http.Error(w, "Formato de archivo no permitido", http.StatusBadRequest)
		return
	}

	// Generar nombre único
	randomName := uuid.New().String() + "." + extension
	objectName := fmt.Sprintf("%s/%s", userData.UserName, randomName)

	// Leemos la imagen en memoria para poder usarla 2 veces
	var buf bytes.Buffer
	size, err := io.Copy(&buf, file)
	if err != nil {
		http.Error(w, "Error al leer el archivo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Obtener dimensiones
	configImage, format, err := image.DecodeConfig(bytes.NewReader(buf.Bytes()))
	if err != nil {
		http.Error(w, "Error al obtener dimensiones de la imagen: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Subir a MinIO
	err = storage.MinioClientInstance.UploadBytes(r.Context(), objectName, buf.Bytes(), "image/"+format)
	if err != nil {
		http.Error(w, "Error al subir imagen a MinIO: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Guardar en la base de datos
	imageData := models.Image{
		Name:     randomName,
		UserName: userData.UserName,
		Path:     objectName,
		Size:     size,
		Format:   format,
		Width:    configImage.Width,
		Height:   configImage.Height,
	}

	err = db.CreateImage(imageData)
	if err != nil {
		http.Error(w, "Error al almacenar la imagen en la base de datos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("%s:%s/%s/%s/%s", config.Cnf.BaseURL, config.Cnf.Port, "images", userData.UserName, randomName)

	resp := models.UploadResponse{
		Message: "Imagen subida exitosamente",
		Image: models.UploadedImageDetail{
			URL:    url,
			Name:   randomName,
			Size:   size,
			Format: format,
			Width:  configImage.Width,
			Height: configImage.Height,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
