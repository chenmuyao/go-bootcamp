package connpool

import (
	"context"
	"database/sql"
	"errors"

	"github.com/chenmuyao/go-bootcamp/pkg/logger"
	"go.uber.org/atomic"
	"gorm.io/gorm"
)

var errUnknownPattern = errors.New("unknown pattern")

const (
	PatternSrcOnly  = "src_only"
	PatternSrcFirst = "src_first"
	PatternDstFirst = "dst_first"
	PatternDstOnly  = "dst_only"
)

type DoubleWriteTx struct {
	src     *sql.Tx
	dst     *sql.Tx
	pattern string
	l       logger.Logger
}

// Commit implements gorm.TxCommitter.
func (d *DoubleWriteTx) Commit() error {
	switch d.pattern {
	case PatternSrcOnly:
		return d.src.Commit()
	case PatternSrcFirst:
		err := d.src.Commit()
		if err != nil {
			return err
		}
		if d.dst != nil {
			er := d.dst.Commit()
			if er != nil {
				d.l.Error(
					"failed to commit transaction for target db in double write",
					logger.Error(err),
				)
			}

		}
		return nil
	case PatternDstFirst:
		err := d.dst.Commit()
		if err != nil {
			return err
		}
		if d.src != nil {
			er := d.src.Commit()
			if er != nil {
				d.l.Error(
					"failed to commit transaction for base db in double write",
					logger.Error(err),
				)
			}

		}
		return nil
	case PatternDstOnly:
		return d.dst.Commit()
	default:
		return errUnknownPattern
	}
}

// Rollback implements gorm.TxCommitter.
func (d *DoubleWriteTx) Rollback() error {
	switch d.pattern {
	case PatternSrcOnly:
		return d.src.Rollback()
	case PatternSrcFirst:
		err := d.src.Rollback()
		if err != nil {
			return err
		}
		if d.dst != nil {
			er := d.dst.Rollback()
			if er != nil {
				d.l.Error(
					"failed to rollback transaction for target db in double write",
					logger.Error(err),
				)
			}

		}
		return nil
	case PatternDstFirst:
		err := d.dst.Rollback()
		if err != nil {
			return err
		}
		if d.src != nil {
			er := d.src.Rollback()
			if er != nil {
				d.l.Error(
					"failed to rollback transaction for base db in double write",
					logger.Error(err),
				)
			}

		}
		return nil
	case PatternDstOnly:
		return d.dst.Rollback()
	default:
		return errUnknownPattern
	}
}

// ExecContext implements gorm.ConnPool.
func (d *DoubleWriteTx) ExecContext(
	ctx context.Context,
	query string,
	args ...interface{},
) (sql.Result, error) {
	switch d.pattern {
	case PatternSrcOnly:
		return d.src.ExecContext(ctx, query, args...)
	case PatternSrcFirst:
		res, err := d.src.ExecContext(ctx, query, args...)
		if err == nil {
			_, er := d.dst.ExecContext(ctx, query, args...)
			if er != nil {
				d.l.Error(
					"failed to write to dst db",
					logger.Error(err),
					logger.String("sql", query),
				)
			}
		}
		return res, err
	case PatternDstFirst:
		res, err := d.dst.ExecContext(ctx, query, args...)
		if err == nil {
			_, er := d.src.ExecContext(ctx, query, args...)
			if er != nil {
				d.l.Error(
					"failed to write to src db",
					logger.Error(err),
					logger.String("sql", query),
				)
			}
		}
		return res, err
	case PatternDstOnly:
		return d.dst.ExecContext(ctx, query, args...)
	default:
		return nil, errUnknownPattern
	}
}

// PrepareContext implements gorm.ConnPool.
func (d *DoubleWriteTx) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	panic("unsupported by double-write")
}

// QueryContext implements gorm.ConnPool.
func (d *DoubleWriteTx) QueryContext(
	ctx context.Context,
	query string,
	args ...interface{},
) (*sql.Rows, error) {
	switch d.pattern {
	case PatternSrcOnly, PatternSrcFirst:
		return d.src.QueryContext(ctx, query, args...)
	case PatternDstOnly, PatternDstFirst:
		return d.dst.QueryContext(ctx, query, args...)
	default:
		return nil, errUnknownPattern
	}
}

// QueryRowContext implements gorm.ConnPool.
func (d *DoubleWriteTx) QueryRowContext(
	ctx context.Context,
	query string,
	args ...interface{},
) *sql.Row {
	switch d.pattern {
	case PatternSrcOnly, PatternSrcFirst:
		return d.src.QueryRowContext(ctx, query, args...)
	case PatternDstOnly, PatternDstFirst:
		return d.dst.QueryRowContext(ctx, query, args...)
	default:
		// XXX: cannot return the error message
		panic(errUnknownPattern)
	}
}

var (
	_ gorm.ConnPool    = &DoubleWriteTx{}
	_ gorm.TxCommitter = &DoubleWriteTx{}
)

type DoubleWritePool struct {
	src     gorm.ConnPool
	dst     gorm.ConnPool
	pattern atomic.String
	l       logger.Logger
}

// BeginTx implements gorm.ConnPoolBeginner.
func (d *DoubleWritePool) BeginTx(ctx context.Context, opts *sql.TxOptions) (gorm.ConnPool, error) {
	pattern := d.pattern.Load()
	switch pattern {
	case PatternSrcOnly:
		src, err := d.src.(gorm.TxBeginner).BeginTx(ctx, opts)
		return &DoubleWriteTx{src: src, l: d.l, pattern: pattern}, err
	case PatternSrcFirst:
		src, err := d.src.(gorm.TxBeginner).BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}
		dst, err := d.dst.(gorm.TxBeginner).BeginTx(ctx, opts)
		if err != nil {
			d.l.Error(
				"Failed to open transaction for the target db in double-write",
				logger.Error(err),
			)
		}
		return &DoubleWriteTx{src: src, dst: dst, l: d.l, pattern: pattern}, nil
	case PatternDstOnly:
		dst, err := d.dst.(gorm.TxBeginner).BeginTx(ctx, opts)
		return &DoubleWriteTx{dst: dst, l: d.l, pattern: pattern}, err
	case PatternDstFirst:
		dst, err := d.dst.(gorm.TxBeginner).BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}
		src, err := d.src.(gorm.TxBeginner).BeginTx(ctx, opts)
		if err != nil {
			d.l.Error(
				"Failed to open transaction for the base db in double-write",
				logger.Error(err),
			)
		}
		return &DoubleWriteTx{src: src, dst: dst, l: d.l, pattern: pattern}, nil
	default:
		return nil, errUnknownPattern
	}
}

func (d *DoubleWritePool) UpdatePattern(pattern string) error {
	switch pattern {
	case PatternSrcOnly, PatternSrcFirst, PatternDstFirst, PatternDstOnly:
		d.pattern.Store(pattern)
		return nil
	default:
		return errUnknownPattern
	}
}

// ExecContext implements gorm.ConnPool.
func (d *DoubleWritePool) ExecContext(
	ctx context.Context,
	query string,
	args ...interface{},
) (sql.Result, error) {
	switch d.pattern.Load() {
	case PatternSrcOnly:
		return d.src.ExecContext(ctx, query, args...)
	case PatternSrcFirst:
		res, err := d.src.ExecContext(ctx, query, args...)
		if err == nil && d.dst != nil {
			_, er := d.dst.ExecContext(ctx, query, args...)
			if er != nil {
				d.l.Error(
					"failed to write to dst db",
					logger.Error(err),
					logger.String("sql", query),
				)
			}
		}
		return res, err
	case PatternDstFirst:
		res, err := d.dst.ExecContext(ctx, query, args...)
		if err == nil && d.src != nil {
			_, er := d.src.ExecContext(ctx, query, args...)
			if er != nil {
				d.l.Error(
					"failed to write to src db",
					logger.Error(err),
					logger.String("sql", query),
				)
			}
		}
		return res, err
	case PatternDstOnly:
		return d.dst.ExecContext(ctx, query, args...)
	default:
		return nil, errUnknownPattern
	}
}

// PrepareContext implements gorm.ConnPool.
func (d *DoubleWritePool) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	panic("unsupported by double-write")
}

// QueryContext implements gorm.ConnPool.
func (d *DoubleWritePool) QueryContext(
	ctx context.Context,
	query string,
	args ...interface{},
) (*sql.Rows, error) {
	switch d.pattern.Load() {
	case PatternSrcOnly, PatternSrcFirst:
		return d.src.QueryContext(ctx, query, args...)
	case PatternDstOnly, PatternDstFirst:
		return d.dst.QueryContext(ctx, query, args...)
	default:
		return nil, errUnknownPattern
	}
}

// QueryRowContext implements gorm.ConnPool.
func (d *DoubleWritePool) QueryRowContext(
	ctx context.Context,
	query string,
	args ...interface{},
) *sql.Row {
	switch d.pattern.Load() {
	case PatternSrcOnly, PatternSrcFirst:
		return d.src.QueryRowContext(ctx, query, args...)
	case PatternDstOnly, PatternDstFirst:
		return d.dst.QueryRowContext(ctx, query, args...)
	default:
		// XXX: cannot return the error message
		panic(errUnknownPattern)
	}
}

func NewDoubleWritePool(
	src *gorm.DB,
	dst *gorm.DB,
	l logger.Logger,
) gorm.ConnPool {
	return &DoubleWritePool{
		src:     src.ConnPool,
		dst:     dst.ConnPool,
		pattern: *atomic.NewString(PatternSrcOnly),
		l:       l,
	}
}

var (
	_ gorm.ConnPool         = &DoubleWritePool{}
	_ gorm.ConnPoolBeginner = &DoubleWritePool{}
)
