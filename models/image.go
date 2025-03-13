package models

import "time"

type Image struct {
	ImageID   int64     `json:"image_id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"not null"`
	UserName  string    `json:"user_name" gorm:"not null;index"`
	User      User      `gorm:"foreignKey:UserName;references:UserName;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Path      string    `json:"path" gorm:"not null"`   // Ruta en el servidor
	Size      int64     `json:"size" gorm:"not null"`   // Tamaño en bytes
	Format    string    `json:"format" gorm:"not null"` // Formato (jpeg, png, etc.)
	Width     int       `json:"width"`                  // Ancho de la imagen
	Height    int       `json:"height"`                 // Alto de la imagen
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
