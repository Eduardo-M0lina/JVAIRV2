package file

import (
	"time"
)

// File representa un archivo asociado a una entidad polimórfica
type File struct {
	ID           int64      `json:"id"`
	Type         *string    `json:"type,omitempty"`
	Path         *string    `json:"path,omitempty"`
	URL          string     `json:"url"`
	FileableID   int64      `json:"fileableId"`
	FileableType string     `json:"fileableType"`
	CreatedAt    *time.Time `json:"createdAt,omitempty"`
	UpdatedAt    *time.Time `json:"updatedAt,omitempty"`
}

// GetDisplayType retorna el tipo de archivo para visualización
func (f *File) GetDisplayType() string {
	if f.Type == nil {
		return "file"
	}
	t := *f.Type
	switch {
	case len(t) > 6 && t[:6] == "image/":
		return "image"
	case len(t) > 6 && t[:6] == "audio/":
		return "audio"
	case len(t) > 6 && t[:6] == "video/":
		return "video"
	default:
		return "file"
	}
}
