package migrations

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type migration1708079654 struct{}

func (m *migration1708079654) GetTime() int32 {
	return 1708079654
}

func (m *migration1708079654) GetName() string {
	return "init table"
}

func (m *migration1708079654) Execute(tx pgx.Tx, ctx context.Context) error {
	queries := []string{
		"CREATE TABLE calculations (id serial PRIMARY KEY , created_at timestamptz, a int, b int, result int, type varchar(255));",
	}
	for _, q := range queries {
		_, err := tx.Exec(ctx, q)
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	addMigration(&migration1708079654{})
}
