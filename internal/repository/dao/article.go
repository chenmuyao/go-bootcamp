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
	GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]Article, error)
	GetByID(ctx context.Context, id int64) (Article, error)
	GetPubByID(ctx context.Context, id int64) (PublishedArticle, error)
}

type GORMArticleDAO struct {
	db *gorm.DB
}

// GetPubByID implements ArticleDAO.
func (a *GORMArticleDAO) GetPubByID(ctx context.Context, id int64) (PublishedArticle, error) {
	var article PublishedArticle
	err := a.db.WithContext(ctx).Where("id = ?", id).First(&article).Error
	return article, err
}

// GetByID implements ArticleDAO.
func (a *GORMArticleDAO) GetByID(ctx context.Context, id int64) (Article, error) {
	var article Article
	err := a.db.WithContext(ctx).Where("id = ?", id).First(&article).Error
	return article, err
}

type Article struct {
	ID      int64  `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"`
	Title   string `gorm:"type=varchar(4096)"       bson:"title,omitempty"`
	Content string `gorm:"type=BLOB"                bson:"content,omitempty"`

	AuthorID int64 `gorm:"index" bson:"author_id,omitempty"`
	Status   uint8 `             bson:"status,omitempty"`
	Ctime    int64 `             bson:"ctime,omitempty"`
	Utime    int64 `             bson:"utime,omitempty"`
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
		if err != nil {
			return 0, err
		}
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
			if err != nil {
				return err
			}
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
	res := a.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":   article.Title,
			"content": article.Content,
			"status":  article.Status,
			"utime":   article.Utime,
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

// GetByAuthor implements ArticleDAO.
func (a *GORMArticleDAO) GetByAuthor(
	ctx context.Context,
	uid int64,
	offset int,
	limit int,
) ([]Article, error) {
	var articles []Article
	err := a.db.WithContext(ctx).
		Where("author_id = ?", uid).
		Offset(offset).
		Limit(limit).
		Order("utime DESC").
		Find(&articles).
		Error
	return articles, err
}

func NewArticleDAO(db *gorm.DB) ArticleDAO {
	return &GORMArticleDAO{
		db: db,
	}
}
