package models

import "time"

type UserModel struct {
	UserName string `gorm:"primaryKey"`
	Password string `gorm:"not null"`
}

func (UserModel) TableName() string {
	return "users"
}

type ImageModel struct {
	ImageID   int64     `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"not null"`
	UserName  string    `gorm:"not null;index"`
	User      UserModel `gorm:"foreignKey:UserName;references:UserName;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Path      string    `gorm:"not null"`
	Size      int64     `gorm:"not null"`
	Format    string    `gorm:"not null"`
	Width     int
	Height    int
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (ImageModel) TableName() string {
	return "images"
}
