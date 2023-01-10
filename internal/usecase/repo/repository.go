package repo

import "context"

type Repository interface {
	Migrate(ctx context.Context) error
}
