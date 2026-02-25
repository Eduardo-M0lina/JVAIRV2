package job_equipment

import (
	"fmt"
	"strings"
	"time"
)

// JobEquipment representa un equipo asociado a un trabajo
type JobEquipment struct {
	ID                  int64      `json:"id"`
	JobID               int64      `json:"jobId"`
	Type                string     `json:"type"`
	Area                *string    `json:"area,omitempty"`
	OutdoorBrand        *string    `json:"outdoorBrand,omitempty"`
	OutdoorModel        *string    `json:"outdoorModel,omitempty"`
	OutdoorSerial       *string    `json:"outdoorSerial,omitempty"`
	OutdoorInstalled    *time.Time `json:"outdoorInstalled,omitempty"`
	FurnaceBrand        *string    `json:"furnaceBrand,omitempty"`
	FurnaceModel        *string    `json:"furnaceModel,omitempty"`
	FurnaceSerial       *string    `json:"furnaceSerial,omitempty"`
	FurnaceInstalled    *time.Time `json:"furnaceInstalled,omitempty"`
	EvaporatorBrand     *string    `json:"evaporatorBrand,omitempty"`
	EvaporatorModel     *string    `json:"evaporatorModel,omitempty"`
	EvaporatorSerial    *string    `json:"evaporatorSerial,omitempty"`
	EvaporatorInstalled *time.Time `json:"evaporatorInstalled,omitempty"`
	AirHandlerBrand     *string    `json:"airHandlerBrand,omitempty"`
	AirHandlerModel     *string    `json:"airHandlerModel,omitempty"`
	AirHandlerSerial    *string    `json:"airHandlerSerial,omitempty"`
	AirHandlerInstalled *time.Time `json:"airHandlerInstalled,omitempty"`
	CreatedAt           *time.Time `json:"createdAt,omitempty"`
	UpdatedAt           *time.Time `json:"updatedAt,omitempty"`
}

// ValidTypes contiene los valores válidos para el campo Type
var ValidTypes = []string{"current", "new"}

// Validate valida los campos requeridos del equipo de trabajo
func (e *JobEquipment) Validate() error {
	if e.JobID == 0 {
		return fmt.Errorf("job_id is required")
	}

	if e.Area == nil || strings.TrimSpace(*e.Area) == "" {
		return fmt.Errorf("area is required")
	}

	if !isValidType(e.Type) {
		return fmt.Errorf("type must be one of: current, new")
	}

	return nil
}

func isValidType(t string) bool {
	for _, v := range ValidTypes {
		if v == t {
			return true
		}
	}
	return false
}

// GetOutdoorUnit retorna la unidad exterior formateada (brand model | S/N serial)
func (e *JobEquipment) GetOutdoorUnit() string {
	parts := []string{}

	if e.OutdoorBrand != nil && *e.OutdoorBrand != "" {
		parts = append(parts, *e.OutdoorBrand)
	}
	if e.OutdoorModel != nil && *e.OutdoorModel != "" {
		parts = append(parts, *e.OutdoorModel)
	}
	if e.OutdoorSerial != nil && *e.OutdoorSerial != "" {
		parts = append(parts, "| S/N "+*e.OutdoorSerial)
	}

	return strings.Join(parts, " ")
}
