package models

type Image struct {
	ImageID  int64  `json:"image_id" gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name" gorm:"not null"`
	UserName string `json:"user_name" gorm:"not null;index"`
	User     User   `gorm:"foreignKey:UserName;references:UserName;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
