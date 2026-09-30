package port

import "context"

// UnitOfWork defines the contract for executing multiple operations within a single transaction.
type UnitOfWork interface {
	// Execute runs the provided function within a transactional context.
	// If the function returns an error, the transaction is automatically rolled back.
	Execute(ctx context.Context, fn func(txCtx context.Context) error) error
}
