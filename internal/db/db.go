package db

import (
	"context"
	"fmt"
	"github.com/MuddCreates/hyperschedule-api-go/internal/data"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UpdateSummary struct {
	Copied   int64
	Upserted int64
	Deleted  int64
}

type UpdateSummaries struct {
	Course  UpdateSummary
	Term    UpdateSummary
	Staff   UpdateSummary
	Section UpdateSummariesSections
}

type Connection struct {
	db *pgxpool.Pool
}

func (conn *Connection) Tx(
	ctx context.Context,
	f func(context.Context, pgx.Tx) error,
) error {
	return conn.db.BeginFunc(ctx, func(tx pgx.Tx) error { return f(ctx, tx) })
}

type Tx struct {
	tx  pgx.Tx
	ctx context.Context
}

func (conn *Connection) Begin(ctx context.Context) (*Tx, error) {
	tx, err := conn.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &Tx{tx, ctx}, nil
}

func (conn *Connection) Batch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return conn.db.SendBatch(ctx, b)
}

func (tx *Tx) Batch(b *pgx.Batch) pgx.BatchResults {
	return tx.tx.SendBatch(tx.ctx, b)
}

func (tx *Tx) Cleanup() error {
	return tx.tx.Rollback(tx.ctx)
}

func (tx *Tx) Commit() error {
	return tx.tx.Commit(tx.ctx)
}

func (tx *Tx) Exec(q string, args ...interface{}) (pgconn.CommandTag, error) {
	return tx.tx.Exec(tx.ctx, q, args...)
}

func (tx *Tx) Query(q string, args ...interface{}) (pgx.Rows, error) {
	return tx.tx.Query(tx.ctx, q, args...)
}

func (tx *Tx) CopyFrom(
	tbl pgx.Identifier, cols []string, rows pgx.CopyFromSource,
) (int64, error) {
	return tx.tx.CopyFrom(tx.ctx, tbl, cols, rows)
}

type handle struct {
	tx  pgx.Tx
	ctx context.Context
}

func New(ctx context.Context, addr string) (*Connection, error) {
	db, err := pgxpool.Connect(ctx, addr)
	if err != nil {
		return nil, err
	}
	return &Connection{db}, nil
}

func (conn *Connection) Update(ctx context.Context, d *data.Data) (*UpdateSummaries, error) {
	tx, err := conn.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	handle := &handle{tx, ctx}

	termInfo, err := handle.updateTerms(d.Terms)
	if err != nil {
		return nil, fmt.Errorf("update terms failed: %w", err)
	}
	courseInfo, err := handle.updateCourses(d.Courses)
	if err != nil {
		return nil, fmt.Errorf("update courses failed: %w", err)
	}
	staffInfo, err := updateStaff(handle, d.Staff)
	if err != nil {
		return nil, fmt.Errorf("update staff failed: %w", err)
	}
	sectionInfo, err := handle.updateSections(d.CourseSections)
	if err != nil {
		return nil, fmt.Errorf("update section failed: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &UpdateSummaries{
		Term:    termInfo,
		Course:  courseInfo,
		Staff:   staffInfo,
		Section: sectionInfo,
	}, nil
}
