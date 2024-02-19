package migrations

import (
	"context"
	"fmt"
	"github.com/dimitryshirokov/simple-app/internal/internal_error"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"sort"
	"strings"
	"time"
)

var migrations = make(map[int32]migration)

func addMigration(m migration) {
	migrations[m.GetTime()] = m
}

type migration interface {
	GetTime() int32
	GetName() string
	Execute(tx pgx.Tx, ctx context.Context) error
}

func NewMigratorRunner(conn *pgxpool.Pool) *MigratorRunner {
	return &MigratorRunner{
		conn:            conn,
		readyMigrations: make([]migration, 0),
	}
}

type MigratorRunner struct {
	conn *pgxpool.Pool

	readyMigrations []migration
}

func (mr *MigratorRunner) CreateMigration() error {
	temp := `package migrations

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type migration%time% struct{}

func (m *migration%time%) GetTime() int32 {
	return %time%
}

func (m *migration%time%) GetName() string {
	// TODO
	return ""
}

func (m *migration%time%) Execute(tx pgx.Tx, ctx context.Context) error {
	// TODO
	return nil
}

func init() {
	addMigration(&migration%time%{})
}
`
	nt := fmt.Sprintf("%d", time.Now().Unix())
	fn := "./internal/migrations/migration_" + nt + ".go"
	f := strings.ReplaceAll(temp, "%time%", nt)
	if _, e := os.Stat(fn); e == nil {
		return fmt.Errorf("migrations with name %s already exists", fn)
	}
	return os.WriteFile(fn, []byte(f), 0777)
}

func (mr *MigratorRunner) Migrate() error {
	mr.createMigrationsOrder()
	err := mr.createTableIfNotExists()
	if err != nil {
		return err
	}
	return mr.executeAll()
}

func (mr *MigratorRunner) getLastMigration() (int32, error) {
	var ver int32
	err := mr.conn.QueryRow(context.Background(), "SELECT CASE WHEN max(version) IS NULL THEN 0 ELSE max(version) END FROM migration_versions;").Scan(&ver)
	if err != nil {
		return 0, internal_error.NewError("can't get migration version", err, nil)
	}
	return ver, nil
}

func (mr *MigratorRunner) createTableIfNotExists() error {
	_, err := mr.conn.Exec(
		context.Background(),
		"CREATE TABLE IF NOT EXISTS migration_versions (id SERIAL NOT NULL PRIMARY KEY, version INTEGER NOT NULL, name varchar(255) NOT NULL, executed_at INTEGER NOT NULL);",
	)
	return err

}

func (mr *MigratorRunner) createMigrationsOrder() {
	keys := make([]int32, 0)
	for k := range migrations {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	for _, key := range keys {
		m := migrations[key]
		mr.readyMigrations = append(mr.readyMigrations, m)
	}
}

func (mr *MigratorRunner) executeAll() error {
	ctx := context.Background()
	lastVersion, err := mr.getLastMigration()
	if err != nil {
		return internal_error.NewError("can't get last migration version", err, nil)
	}
	tx, err := mr.conn.Begin(ctx)
	if err != nil {
		return internal_error.NewError("can't start transaction", err, nil)
	}
	for _, m := range mr.readyMigrations {
		if m.GetTime() > lastVersion {
			e := m.Execute(tx, ctx)
			if e != nil {
				_ = tx.Rollback(ctx)
				return internal_error.NewError("can't execute migration", e, map[string]interface{}{
					"version": m.GetTime(),
				})
			}
			e = mr.saveInfo(tx, m)
			if e != nil {
				_ = tx.Rollback(ctx)
				return internal_error.NewError("can't save migration info", e, map[string]interface{}{
					"version": m.GetTime(),
				})
			}
		}
	}
	return tx.Commit(ctx)
}

func (mr *MigratorRunner) saveInfo(tx pgx.Tx, m migration) error {
	_, err := tx.Exec(
		context.Background(),
		fmt.Sprintf(
			"INSERT INTO migration_versions (version, name, executed_at) VALUES (%d, '%s', %d)",
			m.GetTime(),
			m.GetName(),
			time.Now().Unix(),
		),
	)
	return err
}
