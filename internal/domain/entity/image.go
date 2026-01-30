package entity

import "time"

type Image struct {
	ID        int64
	Name      string
	UserName  string
	Path      string
	Size      int64
	Format    string
	Width     int
	Height    int
	CreatedAt time.Time
}

func NewImage(name, userName, path, format string, size int64, width, height int) *Image {
	return &Image{
		Name:      name,
		UserName:  userName,
		Path:      path,
		Size:      size,
		Format:    format,
		Width:     width,
		Height:    height,
		CreatedAt: time.Now(),
	}
}
