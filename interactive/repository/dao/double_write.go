package dao

import (
	"context"
	"errors"

	"github.com/chenmuyao/go-bootcamp/pkg/logger"
	"go.uber.org/atomic"
)

var errUnknownPattern = errors.New("unknown pattern")

const (
	PatternSrcOnly  = "src_only"
	PatternSrcFirst = "src_first"
	PatternDstFirst = "dst_first"
	PatternDstOnly  = "dst_only"
)

type DoubleWriteDAO struct {
	src     InteractiveDAO
	dst     InteractiveDAO
	pattern atomic.String
	l       logger.Logger
}

func (d *DoubleWriteDAO) UpdatePattern(pattern string) {
	d.pattern.Store(pattern)
}

// IncrReadCnt implements InteractiveDAO.
func (d *DoubleWriteDAO) IncrReadCnt(ctx context.Context, biz string, bizID int64) error {
	pattern := d.pattern.Load()
	switch pattern {
	case PatternSrcOnly:
		return d.src.IncrReadCnt(ctx, biz, bizID)
	case PatternSrcFirst:
		err := d.src.IncrReadCnt(ctx, biz, bizID)
		if err != nil {
			return err
		}
		err = d.dst.IncrReadCnt(ctx, biz, bizID)
		if err != nil {
			d.l.Error(
				"failed to write to dst db",
				logger.Error(err),
				logger.String("biz", biz),
				logger.Int64("bizID", bizID),
			)
			// NOTE: don't return error.
			// if src ok, then on biz level it's ok
		}
		return nil
	case PatternDstFirst:
		err := d.dst.IncrReadCnt(ctx, biz, bizID)
		if err == nil {
			er := d.src.IncrReadCnt(ctx, biz, bizID)
			if er != nil {
				d.l.Error("failed to write to src db", logger.Error(er))
			}
		}
		return err
	case PatternDstOnly:
		return d.dst.IncrReadCnt(ctx, biz, bizID)
	default:
		return errUnknownPattern
	}
}

func (d *DoubleWriteDAO) Get(ctx context.Context, biz string, bizID int64) (Interactive, error) {
	pattern := d.pattern.Load()
	switch pattern {
	case PatternSrcOnly, PatternSrcFirst:
		return d.src.Get(ctx, biz, bizID)
	case PatternDstOnly, PatternDstFirst:
		return d.dst.Get(ctx, biz, bizID)
	default:
		return Interactive{}, errUnknownPattern
	}
}

// NOTE: validate at the same time (not recommended, coupling)
func (d *DoubleWriteDAO) GetV1(ctx context.Context, biz string, bizID int64) (Interactive, error) {
	pattern := d.pattern.Load()
	switch pattern {
	case PatternSrcOnly, PatternSrcFirst:
		intr, err := d.src.Get(ctx, biz, bizID)
		if err != nil {
			go func() {
				intrDst, er := d.dst.Get(ctx, biz, bizID)
				if er != nil {
					if intr == intrDst {
						// log and notify fixer
					}
				}
			}()
		}
		return intr, err
	case PatternDstOnly, PatternDstFirst:
		return d.dst.Get(ctx, biz, bizID)
	default:
		return Interactive{}, errUnknownPattern
	}
}

// BatchIncrReadCnt implements InteractiveDAO.
func (d *DoubleWriteDAO) BatchIncrReadCnt(
	ctx context.Context,
	bizs []string,
	bizIDs []int64,
) error {
	panic("unimplemented")
}

// DeleteCollectionBiz implements InteractiveDAO.
func (d *DoubleWriteDAO) DeleteCollectionBiz(ctx context.Context, cb UserCollectionBiz) error {
	panic("unimplemented")
}

// DeleteLikeInfo implements InteractiveDAO.
func (d *DoubleWriteDAO) DeleteLikeInfo(
	ctx context.Context,
	biz string,
	bizID int64,
	uid int64,
) error {
	panic("unimplemented")
}

// GetAll implements InteractiveDAO.
func (d *DoubleWriteDAO) GetAll(
	ctx context.Context,
	biz string,
	limit int,
	offset int,
) ([]Interactive, error) {
	panic("unimplemented")
}

// GetByIDs implements InteractiveDAO.
func (d *DoubleWriteDAO) GetByIDs(
	ctx context.Context,
	biz string,
	ids []int64,
) ([]Interactive, error) {
	panic("unimplemented")
}

// GetCollectInfo implements InteractiveDAO.
func (d *DoubleWriteDAO) GetCollectInfo(
	ctx context.Context,
	biz string,
	bizID int64,
	uid int64,
) (UserCollectionBiz, error) {
	panic("unimplemented")
}

// GetLikeInfo implements InteractiveDAO.
func (d *DoubleWriteDAO) GetLikeInfo(
	ctx context.Context,
	biz string,
	bizID int64,
	uid int64,
) (UserLikeBiz, error) {
	panic("unimplemented")
}

// InsertCollectionBiz implements InteractiveDAO.
func (d *DoubleWriteDAO) InsertCollectionBiz(ctx context.Context, cb UserCollectionBiz) error {
	panic("unimplemented")
}

// InsertLikeInfo implements InteractiveDAO.
func (d *DoubleWriteDAO) InsertLikeInfo(
	ctx context.Context,
	biz string,
	bizID int64,
	uid int64,
) error {
	panic("unimplemented")
}

// MustBatchGet implements InteractiveDAO.
func (d *DoubleWriteDAO) MustBatchGet(
	ctx context.Context,
	biz string,
	bizIDs []int64,
) ([]Interactive, error) {
	panic("unimplemented")
}

func NewDoubleWriteDAO(src InteractiveDAO, dst InteractiveDAO, l logger.Logger) InteractiveDAO {
	return &DoubleWriteDAO{
		src:     src,
		dst:     dst,
		pattern: *atomic.NewString(PatternSrcOnly),
		l:       l,
	}
}
