package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// UnitOfWorkMock is a mock implementation of the UnitOfWork interface.
type UnitOfWorkMock struct {
	mock.Mock
}

// Execute runs the given function within a transaction context.
func (m *UnitOfWorkMock) Execute(ctx context.Context, fn func(txCtx context.Context) error) error {
	args := m.Called(ctx)

	if fn != nil {
		if err := fn(ctx); err != nil {
			return err
		}
	}

	return args.Error(0)
}
