package job

import "errors"

var (
	// ErrJobNotFound indica que el job no fue encontrado
	ErrJobNotFound = errors.New("job not found")

	// ErrJobDeleted indica que el job está eliminado
	ErrJobDeleted = errors.New("job is deleted")

	// ErrJobAlreadyClosed indica que el job ya está cerrado
	ErrJobAlreadyClosed = errors.New("job is already closed")

	// ErrInvalidJobCategory indica que la categoría de trabajo no es válida
	ErrInvalidJobCategory = errors.New("invalid job category")

	// ErrInvalidJobPriority indica que la prioridad de trabajo no es válida
	ErrInvalidJobPriority = errors.New("invalid job priority")

	// ErrInvalidJobStatus indica que el estado de trabajo no es válido
	ErrInvalidJobStatus = errors.New("invalid job status")

	// ErrInvalidWorkflow indica que el workflow no es válido
	ErrInvalidWorkflow = errors.New("invalid workflow")

	// ErrInvalidProperty indica que la propiedad no es válida
	ErrInvalidProperty = errors.New("invalid property")

	// ErrInvalidUser indica que el usuario no es válido
	ErrInvalidUser = errors.New("invalid user")

	// ErrInvalidTechnicianJobStatus indica que el estado de técnico no es válido
	ErrInvalidTechnicianJobStatus = errors.New("invalid technician job status")

	// ErrWorkflowHasNoStatuses indica que el workflow no tiene statuses configurados
	ErrWorkflowHasNoStatuses = errors.New("workflow has no statuses configured")
)
