package property_equipment

import (
	"fmt"
	"strings"
	"time"
)

// PropertyEquipment representa un equipo asociado a una propiedad
type PropertyEquipment struct {
	ID                  int64      `json:"id"`
	PropertyID          int64      `json:"propertyId"`
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

// Validate valida los campos requeridos del equipo de propiedad
func (e *PropertyEquipment) Validate() error {
	if e.PropertyID == 0 {
		return fmt.Errorf("property_id is required")
	}

	if e.Area == nil || strings.TrimSpace(*e.Area) == "" {
		return fmt.Errorf("area is required")
	}

	return nil
}

// GetOutdoorUnit retorna la unidad exterior formateada (brand model | S/N serial)
func (e *PropertyEquipment) GetOutdoorUnit() string {
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
