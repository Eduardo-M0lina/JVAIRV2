package job_equipment

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/domain/job_equipment"
)

// Handler maneja las peticiones HTTP para equipos de trabajo
type Handler struct {
	useCase *job_equipment.UseCase
}

// NewHandler crea una nueva instancia del handler de equipos de trabajo
func NewHandler(useCase *job_equipment.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

// RegisterRoutes registra las rutas del handler como sub-recurso de jobs
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/jobs/{jobId}/equipment", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

// CreateJobEquipmentRequest representa la solicitud para crear un equipo de trabajo
type CreateJobEquipmentRequest struct {
	Type                string  `json:"type" example:"current"`
	Area                *string `json:"area" example:"Main Floor"`
	OutdoorBrand        *string `json:"outdoorBrand,omitempty" example:"Carrier"`
	OutdoorModel        *string `json:"outdoorModel,omitempty" example:"24ACC636A003"`
	OutdoorSerial       *string `json:"outdoorSerial,omitempty" example:"1234567890"`
	OutdoorInstalled    *string `json:"outdoorInstalled,omitempty" example:"2023-06-15T00:00:00Z"`
	FurnaceBrand        *string `json:"furnaceBrand,omitempty" example:"Lennox"`
	FurnaceModel        *string `json:"furnaceModel,omitempty" example:"SL280UHV"`
	FurnaceSerial       *string `json:"furnaceSerial,omitempty" example:"0987654321"`
	FurnaceInstalled    *string `json:"furnaceInstalled,omitempty" example:"2023-06-15T00:00:00Z"`
	EvaporatorBrand     *string `json:"evaporatorBrand,omitempty" example:"Carrier"`
	EvaporatorModel     *string `json:"evaporatorModel,omitempty" example:"CNPVP3617ALA"`
	EvaporatorSerial    *string `json:"evaporatorSerial,omitempty" example:"1122334455"`
	EvaporatorInstalled *string `json:"evaporatorInstalled,omitempty" example:"2023-06-15T00:00:00Z"`
	AirHandlerBrand     *string `json:"airHandlerBrand,omitempty" example:"Trane"`
	AirHandlerModel     *string `json:"airHandlerModel,omitempty" example:"GAM5A0A36M21SA"`
	AirHandlerSerial    *string `json:"airHandlerSerial,omitempty" example:"5566778899"`
	AirHandlerInstalled *string `json:"airHandlerInstalled,omitempty" example:"2023-06-15T00:00:00Z"`
}

// UpdateJobEquipmentRequest representa la solicitud para actualizar un equipo de trabajo
type UpdateJobEquipmentRequest struct {
	Type                string  `json:"type" example:"current"`
	Area                *string `json:"area" example:"Main Floor"`
	OutdoorBrand        *string `json:"outdoorBrand,omitempty" example:"Carrier"`
	OutdoorModel        *string `json:"outdoorModel,omitempty" example:"24ACC636A003"`
	OutdoorSerial       *string `json:"outdoorSerial,omitempty" example:"1234567890"`
	OutdoorInstalled    *string `json:"outdoorInstalled,omitempty" example:"2023-06-15T00:00:00Z"`
	FurnaceBrand        *string `json:"furnaceBrand,omitempty" example:"Lennox"`
	FurnaceModel        *string `json:"furnaceModel,omitempty" example:"SL280UHV"`
	FurnaceSerial       *string `json:"furnaceSerial,omitempty" example:"0987654321"`
	FurnaceInstalled    *string `json:"furnaceInstalled,omitempty" example:"2023-06-15T00:00:00Z"`
	EvaporatorBrand     *string `json:"evaporatorBrand,omitempty" example:"Carrier"`
	EvaporatorModel     *string `json:"evaporatorModel,omitempty" example:"CNPVP3617ALA"`
	EvaporatorSerial    *string `json:"evaporatorSerial,omitempty" example:"1122334455"`
	EvaporatorInstalled *string `json:"evaporatorInstalled,omitempty" example:"2023-06-15T00:00:00Z"`
	AirHandlerBrand     *string `json:"airHandlerBrand,omitempty" example:"Trane"`
	AirHandlerModel     *string `json:"airHandlerModel,omitempty" example:"GAM5A0A36M21SA"`
	AirHandlerSerial    *string `json:"airHandlerSerial,omitempty" example:"5566778899"`
	AirHandlerInstalled *string `json:"airHandlerInstalled,omitempty" example:"2023-06-15T00:00:00Z"`
}

// JobEquipmentResponse representa la respuesta de un equipo de trabajo
type JobEquipmentResponse struct {
	ID                  int64   `json:"id" example:"1"`
	JobID               int64   `json:"jobId" example:"100"`
	Type                string  `json:"type" example:"current"`
	Area                *string `json:"area,omitempty" example:"Main Floor"`
	OutdoorBrand        *string `json:"outdoorBrand,omitempty" example:"Carrier"`
	OutdoorModel        *string `json:"outdoorModel,omitempty" example:"24ACC636A003"`
	OutdoorSerial       *string `json:"outdoorSerial,omitempty" example:"1234567890"`
	OutdoorInstalled    *string `json:"outdoorInstalled,omitempty" example:"2023-06-15"`
	FurnaceBrand        *string `json:"furnaceBrand,omitempty" example:"Lennox"`
	FurnaceModel        *string `json:"furnaceModel,omitempty" example:"SL280UHV"`
	FurnaceSerial       *string `json:"furnaceSerial,omitempty" example:"0987654321"`
	FurnaceInstalled    *string `json:"furnaceInstalled,omitempty" example:"2023-06-15"`
	EvaporatorBrand     *string `json:"evaporatorBrand,omitempty" example:"Carrier"`
	EvaporatorModel     *string `json:"evaporatorModel,omitempty" example:"CNPVP3617ALA"`
	EvaporatorSerial    *string `json:"evaporatorSerial,omitempty" example:"1122334455"`
	EvaporatorInstalled *string `json:"evaporatorInstalled,omitempty" example:"2023-06-15"`
	AirHandlerBrand     *string `json:"airHandlerBrand,omitempty" example:"Trane"`
	AirHandlerModel     *string `json:"airHandlerModel,omitempty" example:"GAM5A0A36M21SA"`
	AirHandlerSerial    *string `json:"airHandlerSerial,omitempty" example:"5566778899"`
	AirHandlerInstalled *string `json:"airHandlerInstalled,omitempty" example:"2023-06-15"`
	OutdoorUnit         string  `json:"outdoorUnit" example:"Carrier 24ACC636A003 | S/N 1234567890"`
	CreatedAt           string  `json:"createdAt,omitempty" example:"2024-01-15T10:30:00Z"`
	UpdatedAt           string  `json:"updatedAt,omitempty" example:"2024-01-18T14:20:00Z"`
}

func toResponse(e *job_equipment.JobEquipment) JobEquipmentResponse {
	resp := JobEquipmentResponse{
		ID:               e.ID,
		JobID:            e.JobID,
		Type:             e.Type,
		Area:             e.Area,
		OutdoorBrand:     e.OutdoorBrand,
		OutdoorModel:     e.OutdoorModel,
		OutdoorSerial:    e.OutdoorSerial,
		FurnaceBrand:     e.FurnaceBrand,
		FurnaceModel:     e.FurnaceModel,
		FurnaceSerial:    e.FurnaceSerial,
		EvaporatorBrand:  e.EvaporatorBrand,
		EvaporatorModel:  e.EvaporatorModel,
		EvaporatorSerial: e.EvaporatorSerial,
		AirHandlerBrand:  e.AirHandlerBrand,
		AirHandlerModel:  e.AirHandlerModel,
		AirHandlerSerial: e.AirHandlerSerial,
		OutdoorUnit:      e.GetOutdoorUnit(),
	}

	if e.OutdoorInstalled != nil {
		s := e.OutdoorInstalled.Format("2006-01-02")
		resp.OutdoorInstalled = &s
	}
	if e.FurnaceInstalled != nil {
		s := e.FurnaceInstalled.Format("2006-01-02")
		resp.FurnaceInstalled = &s
	}
	if e.EvaporatorInstalled != nil {
		s := e.EvaporatorInstalled.Format("2006-01-02")
		resp.EvaporatorInstalled = &s
	}
	if e.AirHandlerInstalled != nil {
		s := e.AirHandlerInstalled.Format("2006-01-02")
		resp.AirHandlerInstalled = &s
	}
	if e.CreatedAt != nil {
		resp.CreatedAt = e.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	}
	if e.UpdatedAt != nil {
		resp.UpdatedAt = e.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	return resp
}

func parseTimePtr(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02T15:04:05Z07:00", *s)
	if err != nil {
		t, err = time.Parse("2006-01-02", *s)
		if err != nil {
			return nil
		}
	}
	return &t
}
