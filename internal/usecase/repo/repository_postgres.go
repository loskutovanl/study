package repo

import (
	"context"
	"database/sql"
)

type PostgreSQLClassicRepository struct {
	db *sql.DB
}

func NewPostgreSQLClassicRepository(db *sql.DB) *PostgreSQLClassicRepository {
	return &PostgreSQLClassicRepository{
		db: db,
	}
}

func (r *PostgreSQLClassicRepository) Migrate(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS "friends"(
    id SERIAL PRIMARY KEY,
    user1_id INT NOT NULL,
    user2_id INT NOT NULL
	);`

	_, err := r.db.ExecContext(ctx, query)
	return err
}
