package fixer

import (
	"context"

	"github.com/chenmuyao/go-bootcamp/pkg/migrator"
	"github.com/chenmuyao/go-bootcamp/pkg/migrator/events"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Fixer[T migrator.Entity] struct {
	base    *gorm.DB
	target  *gorm.DB
	columns []string
}

func NewFixer[T migrator.Entity](base *gorm.DB, target *gorm.DB) *Fixer[T] {
	return &Fixer[T]{
		base:   base,
		target: target,
	}
}

func (f *Fixer[T]) FixIncorrect(ctx context.Context, evt events.InconsistentEvent) error {
	// BUG: have concurrency issue when doing double-write.
	switch evt.Type {
	case events.InconsistentEventTypeNEQ:
		var t T
		err := f.base.WithContext(ctx).Where("id = ?", evt.ID).First(&t).Error
		switch err {
		case gorm.ErrRecordNotFound:
			return f.target.WithContext(ctx).Model(&t).Delete("id = ?", evt.ID).Error
		case nil:
			return f.target.WithContext(ctx).Updates(&t).Error
		default:
			return nil
		}
	case events.InconsistentEventTypeTargetMissing:
		var t T
		err := f.base.WithContext(ctx).Where("id = ?", evt.ID).First(&t).Error
		switch err {
		case gorm.ErrRecordNotFound:
			return nil
		case nil:
			return f.target.WithContext(ctx).Create(&t).Error
		default:
			return nil
		}
	case events.InconsistentEventTypeBaseMissing:
		var t T
		return f.target.WithContext(ctx).Model(&t).Delete("id = ?", evt.ID).Error
	}
	return nil
}

func (f *Fixer[T]) Fix(ctx context.Context, evt events.InconsistentEvent) error {
	// NOTE: use upsert instead.
	switch evt.Type {
	case events.InconsistentEventTypeNEQ, events.InconsistentEventTypeTargetMissing:
		var t T
		err := f.base.WithContext(ctx).Where("id = ?", evt.ID).First(&t).Error
		switch err {
		case gorm.ErrRecordNotFound:
			return f.target.WithContext(ctx).Model(&t).Delete("id = ?", evt.ID).Error
		case nil:
			// upsert
			return f.target.WithContext(ctx).Clauses(clause.OnConflict{
				DoUpdates: clause.AssignmentColumns(f.columns),
			}).Create(&t).Error
		default:
			return nil
		}
	case events.InconsistentEventTypeBaseMissing:
		var t T
		return f.target.WithContext(ctx).Model(&t).Delete("id = ?", evt.ID).Error
	}
	return nil
}

type OverwriteFixer[T migrator.Entity] struct {
	base    *gorm.DB
	target  *gorm.DB
	columns []string
}

func NewOverwriteFixer[T migrator.Entity](
	base *gorm.DB,
	target *gorm.DB,
) (*OverwriteFixer[T], error) {
	rows, err := base.Order("id").Rows()
	if err != nil {
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	return &OverwriteFixer[T]{
		base:   base,
		target: target,
        columns: columns,
	}, nil
}

func (f *OverwriteFixer[T]) Fix(ctx context.Context, id int64) error {
	// NOTE: Brutal way ...
	var t T
	err := f.base.WithContext(ctx).Where("id = ?", id).First(&t).Error
	switch err {
	case gorm.ErrRecordNotFound:
		return f.target.WithContext(ctx).Model(&t).Delete("id = ?", id).Error
	case nil:
		// upsert
		return f.target.WithContext(ctx).Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns(f.columns),
		}).Create(&t).Error
	default:
		return nil
	}
}
