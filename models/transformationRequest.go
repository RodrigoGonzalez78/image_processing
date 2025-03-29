package models

type TransformationRequest struct {
	Transformations struct {
		Resize struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"resize"`
		Crop struct {
			Width  int `json:"width"`
			Height int `json:"height"`
			X      int `json:"x"`
			Y      int `json:"y"`
		} `json:"crop"`
		Rotate  float64 `json:"rotate"`
		Format  string  `json:"format"`
		Filters struct {
			Grayscale bool `json:"grayscale"`
			Sepia     bool `json:"sepia"`
		} `json:"filters"`
	} `json:"transformations"`
}
