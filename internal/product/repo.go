package product

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/pkg/errors"
)

type PgRepo struct {
	db *sql.DB
}

func NewPgRepo(db *sql.DB) *PgRepo {
	return &PgRepo{db}
}

const createProductsTableQuery = `
create table if not exists products (
	id text not null primary key,
	data jsonb not null,
	last_modified timestamp not null
)`

func (r *PgRepo) Migrate(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, createProductsTableQuery)
	return err
}

const upsertProductQuery = `
insert into products (id, data, last_modified)
values ($1, $2, $3)
on conflict (id) do update
set data = excluded.data,
    last_modified = excluded.last_modified`

func (r *PgRepo) Upsert(ctx context.Context, products []Model) error {
	tx, err := r.db.Begin()
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}

	for _, product := range products {
		b, err := json.Marshal(product.Data)
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, "json marshal")
		}

		_, err = tx.ExecContext(ctx, upsertProductQuery,
			product.ID, string(b), product.LastModified)
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, "upsert product")
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "commit transaction")
	}

	return nil
}
