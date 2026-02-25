package invoice_payment

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) Create(ctx context.Context, payment *InvoicePayment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *MockService) GetByID(ctx context.Context, invoiceID, id int64) (*InvoicePayment, error) {
	args := m.Called(ctx, invoiceID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*InvoicePayment), args.Error(1)
}

func (m *MockService) ListByInvoiceID(ctx context.Context, invoiceID int64, filters map[string]interface{}, page, pageSize int) ([]*InvoicePayment, int, error) {
	args := m.Called(ctx, invoiceID, filters, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*InvoicePayment), args.Int(1), args.Error(2)
}

func (m *MockService) Update(ctx context.Context, payment *InvoicePayment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *MockService) Delete(ctx context.Context, invoiceID, id int64) error {
	args := m.Called(ctx, invoiceID, id)
	return args.Error(0)
}
