package dao

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrArticleNotFound = errors.New("article not found")

type ArticleDAO interface {
	Insert(ctx context.Context, article Article) (int64, error)
	UpdateByID(ctx context.Context, article Article) error
	Sync(ctx context.Context, article Article) (int64, error)
	Transaction(ctx context.Context, fn func(ctx context.Context, tx any) (any, error)) (any, error)
	Upsert(ctx context.Context, article PublishedArticle) error
	// model is like &Article{} or &PublishedArticle{}
	UpdateStatusByID(
		ctx context.Context,
		model any,
		userID int64,
		articleID int64,
		status uint8,
	) error
}

type GORMArticleDAO struct {
	db *gorm.DB
}

func NewArticleDAO(db *gorm.DB) ArticleDAO {
	return &GORMArticleDAO{
		db: db,
	}
}

type Article struct {
	ID      int64  `gorm:"primaryKey,autoIncrement"`
	Title   string `gorm:"type=varchar(4096)"`
	Content string `gorm:"type=BLOB"`

	AuthorID int64 `gorm:"index"`
	Status   uint8
	Ctime    int64
	Utime    int64
}

// same DB, different tables
type PublishedArticle Article

func (a *GORMArticleDAO) Insert(ctx context.Context, article Article) (int64, error) {
	now := time.Now().UnixMilli()
	article.Ctime = now
	article.Utime = now
	err := a.db.WithContext(ctx).Create(&article).Error
	return article.ID, err
}

func (a *GORMArticleDAO) UpdateByID(ctx context.Context, article Article) error {
	now := time.Now().UnixMilli()
	res := a.db.WithContext(ctx).
		Model(&Article{}).
		Where("id = ? AND author_id = ?", article.ID, article.AuthorID).
		Updates(map[string]any{
			"title":   article.Title,
			"content": article.Content,
			"status":  article.Status,
			"utime":   now,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrArticleNotFound
	}
	return nil
}

func (a *GORMArticleDAO) UpdateStatusByID(
	ctx context.Context,
	model any,
	userID int64,
	articleID int64,
	status uint8,
) error {
	now := time.Now().UnixMilli()
	res := a.db.WithContext(ctx).
		Model(&Article{}).
		Where("id = ? AND author_id = ?", articleID, userID).
		Updates(map[string]any{
			"status": status,
			"utime":  now,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrArticleNotFound
	}
	return nil
}

func (a *GORMArticleDAO) SyncV1(ctx context.Context, article Article) (int64, error) {
	tx := a.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	// avoid panicking
	defer tx.Rollback()

	txDAO := NewArticleDAO(tx)
	var err error
	id := article.ID

	if id > 0 {
		err = txDAO.UpdateByID(ctx, article)
	} else {
		id, err = txDAO.Insert(ctx, article)
		if err != nil {
			return 0, err
		}
		article.ID = id
	}
	publishedArticle := PublishedArticle(article)
	err = txDAO.Upsert(ctx, publishedArticle)
	if err != nil {
		return 0, err
	}
	tx.Commit()
	return id, nil
}

// Use closure (recommended)
func (a *GORMArticleDAO) Sync(ctx context.Context, article Article) (int64, error) {
	id := article.ID
	err := a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txDAO := NewArticleDAO(tx)
		var err error
		if id > 0 {
			err = txDAO.UpdateByID(ctx, article)
		} else {
			id, err = txDAO.Insert(ctx, article)
			if err != nil {
				return err
			}
			article.ID = id
		}
		publishedArticle := PublishedArticle(article)
		err = txDAO.Upsert(ctx, publishedArticle)
		if err != nil {
			return err
		}
		return nil
	})
	return id, err
}

func (a *GORMArticleDAO) Transaction(
	ctx context.Context,
	fn func(ctx context.Context, tx any) (any, error),
) (any, error) {
	var ret any
	var err error
	err = a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txDAO := NewArticleDAO(tx)
		ret, err = fn(ctx, txDAO)
		return err
	})
	return ret, err
}

func (a *GORMArticleDAO) Upsert(ctx context.Context, article PublishedArticle) error {
	now := time.Now().UnixMilli()
	res := a.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":   article.Title,
			"content": article.Content,
			"status":  article.Status,
			"utime":   now,
		}),
	}).Create(&article)
	if res.Error != nil {
		// TODO: log and retry
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrArticleNotFound
	}
	return nil
}
