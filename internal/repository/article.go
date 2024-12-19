package repository

import (
	"context"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
	"gorm.io/gorm"
)

var ErrArticleNotFound = dao.ErrArticleNotFound

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
	Sync(ctx context.Context, article domain.Article) (int64, error)
}

type CachedArticleRepository struct {
	// 1 DB 1 table
	dao dao.ArticleDAO

	// V1 no transaction, more suitable for 2 DBs
	readerDAO dao.ArticleReaderDAO
	authorDAO dao.ArticleAuthorDAO

	// V2 repository level transaction
	db *gorm.DB
}

func NewArticleRepository(dao dao.ArticleDAO) ArticleRepository {
	return &CachedArticleRepository{
		dao: dao,
	}
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
	return c.dao.UpdateByID(ctx, c.toEntity(article))
}

func (c *CachedArticleRepository) Create(
	ctx context.Context,
	article domain.Article,
) (int64, error) {
	return c.dao.Insert(ctx, c.toEntity(article))
}

func (c *CachedArticleRepository) Sync(ctx context.Context, article domain.Article) (int64, error) {
	panic("unimplemented")
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

// func (c *CachedArticleRepository) toDAO(article dao.Article) domain.Article {
// 	return domain.Article{
// 		Title:   article.Title,
// 		Content: article.Content,
// 	}
// }

func (c *CachedArticleRepository) toEntity(article domain.Article) dao.Article {
	return dao.Article{
		ID:       article.ID,
		Title:    article.Title,
		Content:  article.Content,
		AuthorID: article.Author.ID,
	}
}
