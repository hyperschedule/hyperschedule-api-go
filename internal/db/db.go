package db

import (
	"context"
	"fmt"
	"github.com/MuddCreates/hyperschedule-api-go/internal/data"
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
