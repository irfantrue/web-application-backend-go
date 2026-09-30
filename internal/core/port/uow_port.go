package port

import "context"

type UnitOfWork interface {
	Execute(ctx context.Context, fn func(txCtx context.Context) error) error
}
