package validator

import (
	"context"
	"time"

	"github.com/chenmuyao/generique/gslice"
	"github.com/chenmuyao/go-bootcamp/pkg/logger"
	"github.com/chenmuyao/go-bootcamp/pkg/migrator"
	"github.com/chenmuyao/go-bootcamp/pkg/migrator/events"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type Validator[T migrator.Entity] struct {
	base      *gorm.DB
	target    *gorm.DB
	l         logger.Logger
	producer  events.Producer
	direction string
	batchSize int
}

func (v *Validator[T]) Validate(ctx context.Context) error {
	var eg errgroup.Group
	eg.Go(func() error {
		return v.ValidateTargetToBase(ctx)
	})
	eg.Go(func() error {
		return v.ValidateBaseToTarget(ctx)
	})
	return eg.Wait()
}

func (v *Validator[T]) ValidateBaseToTarget(ctx context.Context) error {
	offset := -1
	for {
		offset++
		var src T
		err := v.base.WithContext(ctx).Order("id").First(&src).Error
		if err == gorm.ErrRecordNotFound {
			// finished
			return nil
		}
		if err != nil {
			// error
			v.l.Error("base -> target failed to query base", logger.Error(err))
			continue
		}

		var dst T
		err = v.target.WithContext(ctx).Where("id = ?", src.ID()).First(&dst).Error
		switch err {
		case gorm.ErrRecordNotFound:
			// target not found
			v.notify(src.ID(), events.InconsistentEventTypeTargetMissing)
		case nil:
			equal := src.CompareTo(dst)
			if !equal {
				// send a message to kafka
				v.notify(src.ID(), events.InconsistentEventTypeNEQ)
			}
		default:
			v.l.Error(
				"base -> target failed to query target",
				logger.Int64("id", src.ID()),
				logger.Error(err),
			)
		}
	}
}

func (v *Validator[T]) ValidateTargetToBase(ctx context.Context) error {
	offset := -v.batchSize
	for {
		offset += v.batchSize
		var ts []T
		err := v.target.WithContext(ctx).
			Select("id").
			Order("id").
			Offset(offset).
			Limit(v.batchSize).
			Find(&ts).
			Error
		if err == gorm.ErrRecordNotFound || len(ts) == 0 {
			return nil
		}
		if err != nil {
			v.l.Error("target -> base faield to query target", logger.Error(err))
		}
		var srcTs []T
		ids := gslice.Map(ts, func(id int, src T) int64 {
			return src.ID()
		})
		err = v.base.WithContext(ctx).Select("id").Where("id IN ?", ids).Find(&srcTs).Error
		if err == gorm.ErrRecordNotFound || len(ts) == 0 {
			// no data at all
			v.notifyBaseMissing(ts)
			continue
		}
		if err != nil {
			v.l.Error("target -> base failed to query base", logger.Error(err))
			continue
		}

		diff := gslice.DiffSetFunc(srcTs, ts, func(src, dst T) bool {
			return src.ID() == dst.ID()
		})

		v.notifyBaseMissing(diff)

		if len(ts) < v.batchSize {
			return nil
		}
	}
}

func (v *Validator[T]) notifyBaseMissing(ts []T) {
	for _, val := range ts {
		v.notify(val.ID(), events.InconsistentEventTypeBaseMissing)
	}
}

func (v *Validator[T]) notify(id int64, typ string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := v.producer.ProduceInconsistentEvent(
		ctx,
		events.InconsistentEvent{
			ID:        id,
			Direction: v.direction,
			Type:      typ,
		},
	)
	if err != nil {
		v.l.Error(
			"failed to send inconsistent message",
			logger.Error(err),
			logger.Int64("id", id),
			logger.String("type", typ),
			logger.String("direction", v.direction),
		)
	}
}
