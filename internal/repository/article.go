package repository

import (
	"context"
	"time"

	"github.com/chenmuyao/generique/gslice"
	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository/cache"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
	"github.com/chenmuyao/go-bootcamp/pkg/logger"
	"gorm.io/gorm"
)

const pageSize = 100

var ErrArticleNotFound = dao.ErrArticleNotFound

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
	Sync(ctx context.Context, article domain.Article) (int64, error)
	SyncStatus(
		ctx context.Context,
		userID int64,
		articleID int64,
		status domain.ArticleStatus,
	) error
	GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
}

type CachedArticleRepository struct {
	l     logger.Logger
	cache cache.ArticleCache
	// 1 DB 1 table
	dao dao.ArticleDAO

	// V1 no transaction, more suitable for 2 DBs
	readerDAO dao.ArticleReaderDAO
	authorDAO dao.ArticleAuthorDAO

	// V2 repository level transaction
	db *gorm.DB
}

func NewArticleRepositoryV2(
	readerDAO dao.ArticleReaderDAO,
	authorDAO dao.ArticleAuthorDAO,
) *CachedArticleRepository {
	return &CachedArticleRepository{
		readerDAO: readerDAO,
		authorDAO: authorDAO,
	}
}

func (c *CachedArticleRepository) Update(ctx context.Context, article domain.Article) error {
	err := c.dao.UpdateByID(ctx, c.toEntity(article))
	if err != nil {
		return err
	}
	errCache := c.cache.DelFirstPage(ctx, article.Author.ID)
	if errCache != nil {
		c.l.Warn("delete first page cache error", logger.Error(errCache))
	}
	return nil
}

func (c *CachedArticleRepository) Create(
	ctx context.Context,
	article domain.Article,
) (int64, error) {
	id, err := c.dao.Insert(ctx, c.toEntity(article))
	if err != nil {
		return 0, err
	}
	errCache := c.cache.DelFirstPage(ctx, article.Author.ID)
	if errCache != nil {
		c.l.Warn("delete first page cache error", logger.Error(errCache))
	}
	return id, nil
}

func (c *CachedArticleRepository) SyncV0(
	ctx context.Context,
	article domain.Article,
) (int64, error) {
	// Distribute at DAO level
	return c.dao.Sync(ctx, c.toEntity(article))
}

func (c *CachedArticleRepository) Sync(
	ctx context.Context,
	article domain.Article,
) (int64, error) {
	ret, err := c.dao.Transaction(ctx, func(ctx context.Context, tx any) (any, error) {
		daoTx := tx.(dao.ArticleDAO)
		var err error
		id := article.ID

		articleEntity := c.toEntity(article)
		if id > 0 {
			err = daoTx.UpdateByID(ctx, articleEntity)
			if err != nil {
				return 0, err
			}
		} else {
			id, err = daoTx.Insert(ctx, articleEntity)
			if err != nil {
				return 0, err
			}
			articleEntity.ID = id
		}
		now := time.Now().UnixMilli()
		publishedArticle := dao.PublishedArticle(articleEntity)
		publishedArticle.Utime = now
		publishedArticle.Ctime = now
		err = daoTx.Upsert(ctx, publishedArticle)
		if err != nil {
			return 0, err
		}
		return id, err
	})
	if err != nil {
		// TODO: log and retry
		return 0, err
	}
	id := ret.(int64)
	errCache := c.cache.DelFirstPage(ctx, article.Author.ID)
	if errCache != nil {
		c.l.Warn("delete first page cache error", logger.Error(errCache))
	}
	return id, nil
}

// NOTE: manual transaction --> not recommended
// defauts: Depends on DAO's dependencies. Not interface-oriented.
// WARN: Should avoid as much as possible!
func (c *CachedArticleRepository) SyncV2(
	ctx context.Context,
	article domain.Article,
) (int64, error) {
	tx := c.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	// avoid panicking
	defer tx.Rollback()

	authorDAO := dao.NewArticleGORMAuthorDAO(tx)
	readerDAO := dao.NewArticleGORMReaderDAO(tx)
	var err error
	id := article.ID

	articleEntity := c.toEntity(article)
	if id > 0 {
		err = authorDAO.UpdateByID(ctx, articleEntity)
	} else {
		id, err = authorDAO.Create(ctx, articleEntity)
		if err != nil {
			return 0, err
		}
		articleEntity.ID = id
	}
	err = readerDAO.UpsertV2(ctx, dao.PublishedArticle(articleEntity))
	if err != nil {
		// TODO: log and retry
		return 0, err
	}
	tx.Commit()
	return id, nil
}

func (c *CachedArticleRepository) SyncV1(
	ctx context.Context,
	article domain.Article,
) (int64, error) {
	var err error
	id := article.ID

	articleEntity := c.toEntity(article)
	if id > 0 {
		err = c.authorDAO.UpdateByID(ctx, articleEntity)
	} else {
		id, err = c.authorDAO.Create(ctx, articleEntity)
		if err != nil {
			return 0, err
		}
		articleEntity.ID = id
	}
	err = c.readerDAO.Upsert(ctx, articleEntity)
	if err != nil {
		// TODO: log and retry
		return 0, err
	}
	return id, nil
}

func (c *CachedArticleRepository) SyncStatus(
	ctx context.Context,
	userID int64,
	articleID int64,
	status domain.ArticleStatus,
) error {
	_, err := c.dao.Transaction(ctx, func(ctx context.Context, tx any) (any, error) {
		daoTx := tx.(dao.ArticleDAO)
		err := daoTx.UpdateStatusByID(ctx, &dao.Article{}, userID, articleID, uint8(status))
		if err != nil {
			return nil, err
		}

		err = daoTx.UpdateStatusByID(ctx, &dao.PublishedArticle{}, userID, articleID, uint8(status))
		if err != nil {
			return nil, err
		}

		return nil, nil
	})
	if err != nil {
		// TODO: log and retry
		return err
	}
	errCache := c.cache.DelFirstPage(ctx, userID)
	if errCache != nil {
		c.l.Warn("delete first page cache error", logger.Error(errCache))
	}
	return nil
}

func (c *CachedArticleRepository) toDomain(article dao.Article) domain.Article {
	return domain.Article{
		ID:      article.ID,
		Title:   article.Title,
		Content: article.Content,
		Author: domain.Author{
			ID: article.AuthorID,
		},
		Status: domain.ArticleStatus(article.Status),
		Ctime:  time.UnixMilli(article.Ctime),
		Utime:  time.UnixMilli(article.Utime),
	}
}

func (c *CachedArticleRepository) toEntity(article domain.Article) dao.Article {
	return dao.Article{
		ID:       article.ID,
		Title:    article.Title,
		Content:  article.Content,
		AuthorID: article.Author.ID,
		Status:   uint8(article.Status),
	}
}

// GetByAuthor implements ArticleRepository.
func (c *CachedArticleRepository) GetByAuthor(
	ctx context.Context,
	uid int64,
	offset int,
	limit int,
) ([]domain.Article, error) {
	// check if query the cache
	if offset == 0 && limit <= pageSize {
		cachedArticles, err := c.cache.GetFirstPage(ctx, uid)
		switch err {
		case nil:
			return cachedArticles[:limit], nil
		case cache.ErrKeyNotExist:
			// ignore
			break
		default:
			// log error and monitor
			c.l.Warn("article cache get error", logger.Error(err))
		}
	}
	articlesDAO, err := c.dao.GetByAuthor(ctx, uid, offset, limit)
	if err != nil {
		return nil, err
	}
	res := gslice.Map(
		articlesDAO,
		func(id int, src dao.Article) domain.Article { return c.toDomain(src) },
	)

	go func() {
		if offset == 0 && limit <= pageSize {
			err = c.cache.SetFirstPage(ctx, uid, res)
			// err 1: network issue, maybe temporary
			// err 2: redis down
			if err != nil {
				// WARN: Monitor here
				c.l.Warn("article cache set error", logger.Error(err))
			}
		}
	}()

	return res, nil
}

func NewArticleRepository(
	l logger.Logger,
	dao dao.ArticleDAO,
	cache cache.ArticleCache,
) ArticleRepository {
	return &CachedArticleRepository{
		l:     l,
		dao:   dao,
		cache: cache,
	}
}
