package service

import (
	"context"
	"github.com/dimitryshirokov/simple-app/internal/config"
	"github.com/dimitryshirokov/simple-app/internal/internal_error"
	"github.com/dimitryshirokov/simple-app/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func NewCalculatorService(ctx context.Context, conf *config.Config, dbPool *pgxpool.Pool) *CalculatorService {
	return &CalculatorService{ctx: ctx, conf: conf, dbPool: dbPool}
}

type CalculatorService struct {
	ctx    context.Context
	conf   *config.Config
	dbPool *pgxpool.Pool
}

func (s *CalculatorService) Addition(a, b int) (*model.Calculation, error) {
	c := &model.Calculation{
		CreatedAt: time.Now(),
		A:         a,
		B:         b,
		Result:    a + b,
		Type:      "addition",
	}
	err := s.saveCalculation(c)
	if err != nil {
		return nil, internal_error.NewError("can't process addition", err, map[string]interface{}{
			"a": a,
			"b": b,
		})
	}
	return c, nil
}

func (s *CalculatorService) Subtraction(a, b int) (*model.Calculation, error) {
	c := &model.Calculation{
		CreatedAt: time.Now(),
		A:         a,
		B:         b,
		Result:    a - b,
		Type:      "subtraction",
	}
	err := s.saveCalculation(c)
	if err != nil {
		return nil, internal_error.NewError("can't process subtraction", err, map[string]interface{}{
			"a": a,
			"b": b,
		})
	}
	return c, nil
}

func (s *CalculatorService) saveCalculation(c *model.Calculation) error {
	ctx, cancel := context.WithTimeout(s.ctx, time.Duration(s.conf.QueryTimeout)*time.Second)
	defer cancel()
	var id int
	err := s.dbPool.QueryRow(
		ctx,
		"INSERT INTO calculations (created_at, a, b, result, type) VALUES (@createdAt, @a, @b, @result, @type) RETURNING id",
		pgx.NamedArgs{
			"createdAt": c.CreatedAt,
			"a":         c.A,
			"b":         c.B,
			"result":    c.Result,
			"type":      c.Type,
		},
	).Scan(&id)
	if err != nil {
		return internal_error.NewError("can't insert calculation", err, map[string]interface{}{
			"createdAt": c.CreatedAt,
			"a":         c.A,
			"b":         c.B,
			"result":    c.Result,
			"type":      c.Type,
		})
	}
	c.Id = id
	return nil
}

func (s *CalculatorService) Results(calculationType string, limit int, offset int) ([]*model.Calculation, int, error) {
	ctx, cancel := context.WithTimeout(s.ctx, time.Duration(s.conf.QueryTimeout)*time.Second)
	defer cancel()
	rows, err := s.dbPool.Query(
		ctx,
		"SELECT id, created_at, a, b, result, type FROM calculations WHERE type = @type ORDER BY id DESC LIMIT @limit OFFSET @offset;",
		pgx.NamedArgs{
			"type":   calculationType,
			"limit":  limit,
			"offset": offset,
		},
	)
	if err != nil {
		return nil, 0, internal_error.NewError("can't execute select", err, nil)
	}
	result := make([]*model.Calculation, 0)
	for rows.Next() {
		c := &model.Calculation{}
		err := rows.Scan(&c.Id, &c.CreatedAt, &c.A, &c.B, &c.Result, &c.Type)
		if err != nil {
			return nil, 0, internal_error.NewError("scan error", err, nil)
		}
		result = append(result, c)
	}
	var count int
	err = s.dbPool.QueryRow(ctx, "SELECT count(id) FROM calculations WHERE type = @type", pgx.NamedArgs{
		"type": calculationType,
	}).Scan(&count)
	if err != nil {
		return nil, 0, internal_error.NewError("can't get count of calculations", err, nil)
	}
	return result, count, nil
}
