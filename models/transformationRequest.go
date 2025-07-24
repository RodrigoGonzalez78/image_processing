package models

// TransformationRequest representa las transformaciones que se pueden aplicar a una imagen.
// swagger:model TransformationRequest
type TransformationRequest struct {
	// Transformaciones a aplicar a la imagen
	Transformations struct {
		// Redimensionar la imagen
		Resize struct {
			// Ancho en píxeles
			// example: 800
			Width int `json:"width"`
			// Alto en píxeles
			// example: 600
			Height int `json:"height"`
		} `json:"resize"`

		// Recortar la imagen
		Crop struct {
			// Ancho del recorte
			// example: 400
			Width int `json:"width"`
			// Alto del recorte
			// example: 300
			Height int `json:"height"`
			// Posición X del recorte
			// example: 100
			X int `json:"x"`
			// Posición Y del recorte
			// example: 50
			Y int `json:"y"`
		} `json:"crop"`

		// Rotar la imagen (en grados)
		// example: 90
		Rotate float64 `json:"rotate"`

		// Formato de salida de la imagen (png, jpg, gif)
		// example: jpg
		Format string `json:"format"`

		// Filtros a aplicar
		Filters struct {
			// Convertir a escala de grises
			// example: true
			Grayscale bool `json:"grayscale"`
			// Aplicar filtro sepia
			// example: false
			Sepia bool `json:"sepia"`
		} `json:"filters"`
	} `json:"transformations"`
}
