package email

import (
	"context"
	"fmt"

	domainCustomer "github.com/your-org/jvairv2/pkg/domain/customer"
	domainInvoice "github.com/your-org/jvairv2/pkg/domain/invoice"
	domainJob "github.com/your-org/jvairv2/pkg/domain/job"
	domainJobResident "github.com/your-org/jvairv2/pkg/domain/job_resident"
	domainJobTask "github.com/your-org/jvairv2/pkg/domain/job_task"
	domainProperty "github.com/your-org/jvairv2/pkg/domain/property"
	domainQuote "github.com/your-org/jvairv2/pkg/domain/quote"
	domainUser "github.com/your-org/jvairv2/pkg/domain/user"
)

// JobRepositoryAdapter adapta el repositorio de job
type JobRepositoryAdapter struct {
	repo domainJob.Repository
}

func NewJobRepositoryAdapter(repo domainJob.Repository) JobRepository {
	return &JobRepositoryAdapter{repo: repo}
}

func (a *JobRepositoryAdapter) GetByID(ctx context.Context, id int64) (*JobData, error) {
	job, err := a.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var dispatchDate *string
	if job.DispatchDate != nil {
		d := job.DispatchDate.Format("January 2, 2006")
		dispatchDate = &d
	}

	return &JobData{
		ID:            job.ID,
		WorkOrder:     job.WorkOrder,
		DispatchDate:  dispatchDate,
		DispatchNotes: job.DispatchNotes,
		PropertyID:    job.PropertyID,
		UserID:        job.UserID,
	}, nil
}

// PropertyRepositoryAdapter adapta el repositorio de property
type PropertyRepositoryAdapter struct {
	repo domainProperty.Repository
}

func NewPropertyRepositoryAdapter(repo domainProperty.Repository) PropertyRepository {
	return &PropertyRepositoryAdapter{repo: repo}
}

func (a *PropertyRepositoryAdapter) GetByID(ctx context.Context, id int64) (*PropertyData, error) {
	prop, err := a.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &PropertyData{
		ID:         prop.ID,
		Street:     prop.Street,
		City:       prop.City,
		State:      prop.State,
		Zip:        prop.Zip,
		CustomerID: prop.CustomerID,
	}, nil
}

// CustomerRepositoryAdapter adapta el repositorio de customer
type CustomerRepositoryAdapter struct {
	repo domainCustomer.Repository
}

func NewCustomerRepositoryAdapter(repo domainCustomer.Repository) CustomerRepository {
	return &CustomerRepositoryAdapter{repo: repo}
}

func (a *CustomerRepositoryAdapter) GetByID(ctx context.Context, id int64) (*CustomerData, error) {
	customer, err := a.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &CustomerData{
		ID:   customer.ID,
		Name: customer.Name,
	}, nil
}

// UserRepositoryAdapter adapta el repositorio de user
type UserRepositoryAdapter struct {
	repo domainUser.Repository
}

func NewUserRepositoryAdapter(repo domainUser.Repository) UserRepository {
	return &UserRepositoryAdapter{repo: repo}
}

func (a *UserRepositoryAdapter) GetByID(ctx context.Context, id int64) (*UserData, error) {
	// Convertir int64 a string para el repositorio de user
	userID := fmt.Sprintf("%d", id)
	user, err := a.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &UserData{
		ID:   user.ID,
		Name: user.Name,
	}, nil
}

// ResidentRepositoryAdapter adapta el repositorio de job_resident
type ResidentRepositoryAdapter struct {
	repo domainJobResident.Repository
}

func NewResidentRepositoryAdapter(repo domainJobResident.Repository) ResidentRepository {
	return &ResidentRepositoryAdapter{repo: repo}
}

func (a *ResidentRepositoryAdapter) ListByJobID(ctx context.Context, jobID int64) ([]*ResidentData, error) {
	// Usar List con jobID, limit grande y offset 0
	residents, _, err := a.repo.List(ctx, jobID, 100, 0)
	if err != nil {
		return nil, err
	}

	var result []*ResidentData
	for _, r := range residents {
		result = append(result, &ResidentData{
			Name:        r.Name,
			MobilePhone: r.MobilePhone,
			HomePhone:   r.HomePhone,
		})
	}

	return result, nil
}

// InvoiceRepositoryAdapter adapta el repositorio de invoice
type InvoiceRepositoryAdapter struct {
	repo domainInvoice.Repository
}

func NewInvoiceRepositoryAdapter(repo domainInvoice.Repository) InvoiceRepository {
	return &InvoiceRepositoryAdapter{repo: repo}
}

func (a *InvoiceRepositoryAdapter) GetByID(ctx context.Context, id int64) (*InvoiceData, error) {
	invoice, err := a.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &InvoiceData{
		ID:            invoice.ID,
		InvoiceNumber: invoice.InvoiceNumber,
		Amount:        invoice.Total,
		Description:   invoice.Description,
		Notes:         invoice.Notes,
		JobID:         invoice.JobID,
	}, nil
}

// QuoteRepositoryAdapter adapta el repositorio de quote
type QuoteRepositoryAdapter struct {
	repo domainQuote.Repository
}

func NewQuoteRepositoryAdapter(repo domainQuote.Repository) QuoteRepository {
	return &QuoteRepositoryAdapter{repo: repo}
}

func (a *QuoteRepositoryAdapter) GetByID(ctx context.Context, id int64) (*QuoteData, error) {
	quote, err := a.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &QuoteData{
		ID:          quote.ID,
		QuoteNumber: quote.QuoteNumber,
		Amount:      quote.Amount,
		Description: quote.Description,
		Notes:       quote.Notes,
		Status:      fmt.Sprintf("%d", quote.QuoteStatusID),
		JobID:       quote.JobID,
	}, nil
}

// TaskRepositoryAdapter adapta el repositorio de job_task
type TaskRepositoryAdapter struct {
	repo domainJobTask.Repository
}

func NewTaskRepositoryAdapter(repo domainJobTask.Repository) TaskRepository {
	return &TaskRepositoryAdapter{repo: repo}
}

func (a *TaskRepositoryAdapter) GetByID(ctx context.Context, id int64) (*TaskData, error) {
	task, err := a.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var dueDate *string
	if task.DueDate != nil {
		d := task.DueDate.Format("January 2, 2006")
		dueDate = &d
	}

	return &TaskData{
		ID:          task.ID,
		Description: task.Task,
		DueDate:     dueDate,
		Status:      fmt.Sprintf("%d", task.TaskStatusID),
		Notes:       nil,
		JobID:       task.JobID,
		UserID:      task.UserID,
	}, nil
}
