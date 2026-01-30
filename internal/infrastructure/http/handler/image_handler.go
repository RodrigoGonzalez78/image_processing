package handler

import (
	"bytes"
	"encoding/json"
	"image"
	"io"
	"net/http"
	"path"
	"strconv"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	imageUC "github.com/RodrigoGonzalez78/internal/application/usecase/image"
	"github.com/RodrigoGonzalez78/internal/infrastructure/http/dto"
	"github.com/RodrigoGonzalez78/internal/infrastructure/http/middleware"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ImageHandler struct {
	uploadUC    *imageUC.UploadUseCase
	getUC       *imageUC.GetUseCase
	listUC      *imageUC.ListUseCase
	transformUC *imageUC.TransformUseCase
}

func NewImageHandler(
	uploadUC *imageUC.UploadUseCase,
	getUC *imageUC.GetUseCase,
	listUC *imageUC.ListUseCase,
	transformUC *imageUC.TransformUseCase,
) *ImageHandler {
	return &ImageHandler{
		uploadUC:    uploadUC,
		getUC:       getUC,
		listUC:      listUC,
		transformUC: transformUC,
	}
}

// Upload godoc
// @Summary      Sube una imagen
// @Description  Permite a un usuario autenticado subir una imagen. Se guarda en MinIO y se almacena la metadata.
// @Tags         images
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        image formData file true "Imagen a subir (jpg, jpeg, png, gif)"
// @Success      201 {object} dto.UploadResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /upload [post]
func (h *ImageHandler) Upload(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error al obtener el archivo: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	userData, ok := middleware.GetUserFromContext(r.Context())
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

	// Leer archivo en memoria
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, file); err != nil {
		http.Error(w, "Error al leer el archivo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Obtener dimensiones
	configImage, format, err := image.DecodeConfig(bytes.NewReader(buf.Bytes()))
	if err != nil {
		http.Error(w, "Error al obtener dimensiones de la imagen: "+err.Error(), http.StatusInternalServerError)
		return
	}

	randomName := uuid.New().String() + "." + extension

	input := imageUC.UploadInput{
		FileName:    randomName,
		UserName:    userData.UserName,
		Data:        buf.Bytes(),
		ContentType: "image/" + format,
		Format:      format,
		Width:       configImage.Width,
		Height:      configImage.Height,
	}

	output, err := h.uploadUC.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, "Error al subir imagen: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.UploadResponse{
		Message: "Imagen subida exitosamente",
		Image: dto.UploadedImageDetail{
			URL:    output.URL,
			Name:   output.Name,
			Size:   output.Size,
			Format: output.Format,
			Width:  output.Width,
			Height: output.Height,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GetImage godoc
// @Summary      Obtener información de una imagen
// @Description  Devuelve la metadata y URL de una imagen del usuario autenticado
// @Tags         images
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "ID de la imagen"
// @Success      200 {object} dto.ImageDetailResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Router       /images/{id} [get]
func (h *ImageHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageIDStr := vars["id"]
	imageID, err := strconv.ParseInt(imageIDStr, 10, 64)
	if err != nil {
		http.Error(w, "ID de imagen inválido", http.StatusBadRequest)
		return
	}

	userData, ok := middleware.GetUserFromContext(r.Context())
	if !ok || userData == nil {
		http.Error(w, "No autorizado", http.StatusUnauthorized)
		return
	}

	input := imageUC.GetInput{
		ImageID:  imageID,
		UserName: userData.UserName,
	}

	output, err := h.getUC.Execute(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp := dto.ImageDetailResponse{
		URL:    output.URL,
		Name:   output.Name,
		Size:   output.Size,
		Format: output.Format,
		Width:  output.Width,
		Height: output.Height,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// ListUserImages godoc
// @Summary      Lista imágenes del usuario autenticado
// @Description  Devuelve las imágenes subidas por el usuario autenticado con soporte de paginación.
// @Tags         images
// @Security     BearerAuth
// @Produce      json
// @Param        page   query int false "Número de página" default(1)
// @Param        limit  query int false "Cantidad por página" default(10)
// @Success      200 {object} dto.PaginatedImagesResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /user-images [get]
func (h *ImageHandler) ListUserImages(w http.ResponseWriter, r *http.Request) {
	userData, ok := middleware.GetUserFromContext(r.Context())
	if !ok || userData == nil {
		http.Error(w, "Usuario no autenticado", http.StatusUnauthorized)
		return
	}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	input := imageUC.ListInput{
		UserName: userData.UserName,
		Page:     page,
		Limit:    limit,
	}

	output, err := h.listUC.Execute(input)
	if err != nil {
		http.Error(w, "Error al obtener imágenes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convertir a DTO
	images := make([]dto.ImageItem, len(output.Images))
	for i, img := range output.Images {
		images[i] = dto.ImageItem{
			ID:     img.ID,
			URL:    img.URL,
			Name:   img.Name,
			Size:   img.Size,
			Format: img.Format,
			Width:  img.Width,
			Height: img.Height,
		}
	}

	resp := dto.PaginatedImagesResponse{
		Page:   output.Page,
		Limit:  output.Limit,
		Total:  output.Total,
		Images: images,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// TransformImage godoc
// @Summary      Aplica transformaciones a una imagen
// @Description  Aplica transformaciones (resize, crop, rotación, filtros) a una imagen previamente cargada por el usuario
// @Tags         images
// @Security     BearerAuth
// @Produce      image/png
// @Produce      image/jpeg
// @Produce      image/gif
// @Param        id path int true "ID de la imagen"
// @Param        body body dto.TransformationRequest true "Parámetros de transformación"
// @Success      200 {file} file "Imagen transformada"
// @Failure      400 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /images/{id}/transform [post]
func (h *ImageHandler) TransformImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageIDStr := vars["id"]
	imageID, err := strconv.ParseInt(imageIDStr, 10, 64)
	if err != nil {
		http.Error(w, "ID de imagen inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	userData, ok := middleware.GetUserFromContext(r.Context())
	if !ok || userData == nil {
		http.Error(w, "No autorizado", http.StatusUnauthorized)
		return
	}

	var req dto.TransformationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al parsear JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	input := imageUC.TransformInput{
		ImageID:  imageID,
		UserName: userData.UserName,
		Resize: imageUC.ResizeParams{
			Width:  req.Transformations.Resize.Width,
			Height: req.Transformations.Resize.Height,
		},
		Crop: imageUC.CropParams{
			Width:  req.Transformations.Crop.Width,
			Height: req.Transformations.Crop.Height,
			X:      req.Transformations.Crop.X,
			Y:      req.Transformations.Crop.Y,
		},
		Rotate: req.Transformations.Rotate,
		Format: req.Transformations.Format,
		Filters: imageUC.FilterParams{
			Grayscale: req.Transformations.Filters.Grayscale,
			Sepia:     req.Transformations.Filters.Sepia,
		},
	}

	output, err := h.transformUC.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", output.ContentType)
	w.WriteHeader(http.StatusOK)
	if err := imaging.Encode(w, output.Image, output.Format); err != nil {
		http.Error(w, "Error al codificar la imagen: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

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
func (h *ImageHandler) ServeImage(w http.ResponseWriter, r *http.Request) {
	objectPath := strings.TrimPrefix(r.URL.Path, "/images/")
	objectPath = path.Clean(objectPath)

	if objectPath == "." || strings.Contains(objectPath, "..") {
		http.Error(w, "Acceso prohibido", http.StatusForbidden)
		return
	}

	reader, err := h.transformUC.ServeImage(r.Context(), objectPath)
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
