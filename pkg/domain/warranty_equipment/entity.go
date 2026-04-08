package warranty_equipment

import (
	"fmt"
	"strings"
	"time"
)

// WarrantyEquipment representa un equipo asociado a una garantía
type WarrantyEquipment struct {
	ID                  int64      `json:"id"`
	WarrantyID          int64      `json:"warrantyId"`
	Area                string     `json:"area"`
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

// Validate valida los campos requeridos del equipo de garantía
func (we *WarrantyEquipment) Validate() error {
	if we.WarrantyID == 0 {
		return fmt.Errorf("warranty_id is required")
	}
	if strings.TrimSpace(we.Area) == "" {
		return fmt.Errorf("area is required")
	}
	return nil
}
