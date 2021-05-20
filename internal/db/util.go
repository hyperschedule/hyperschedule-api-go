package db

import (
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type batchOpts struct {
	tmpName         string
	tmpCols         []string
	sqlTmpDefs      string
	sqlInsertSelect string
	rows            [][]interface{}
}

func batchPrepare(h *handle, opts *batchOpts) error {
	stmtCreateTmp := fmt.Sprintf(
		"CREATE TEMPORARY TABLE %s (%s) ON COMMIT DROP",
		pgx.Identifier{opts.tmpName}.Sanitize(),
		opts.sqlTmpDefs,
	)

	if _, err := h.tx.Exec(h.ctx, stmtCreateTmp); err != nil {
		return err
	}
	if _, err := h.tx.CopyFrom(
		h.ctx, pgx.Identifier{opts.tmpName}, opts.tmpCols,
		pgx.CopyFromRows(opts.rows),
	); err != nil {
		return err
	}

	return nil
}

func batchInsert(h *handle, opts *batchOpts) (pgconn.CommandTag, error) {
	if err := batchPrepare(h, opts); err != nil {
		return nil, err
	}
	return h.tx.Exec(h.ctx, opts.sqlInsertSelect)
}

func batchInsertReturns(h *handle, opts *batchOpts) (pgx.Rows, error) {
	if err := batchPrepare(h, opts); err != nil {
		return nil, err
	}
	return h.tx.Query(h.ctx, opts.sqlInsertSelect)
}
