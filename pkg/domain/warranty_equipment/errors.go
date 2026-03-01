package warranty_equipment

import "errors"

var (
	// ErrWarrantyEquipmentNotFound indica que el equipo de garantía no fue encontrado
	ErrWarrantyEquipmentNotFound = errors.New("warranty equipment not found")

	// ErrInvalidWarranty indica que la garantía no es válida
	ErrInvalidWarranty = errors.New("invalid warranty")
)
