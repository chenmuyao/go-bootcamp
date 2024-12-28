package dao

import (
	"context"
	"time"

	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	ArticleCollName          = "articles"
	PublishedArticleCollName = "published_articles"
	DatabaseName             = "wetravel"
)

type MongoDBArticleDAO struct {
	node     *snowflake.Node
	coll     *mongo.Collection
	liveColl *mongo.Collection
	client   *mongo.Client
}

// GetByID implements ArticleDAO.
func (m *MongoDBArticleDAO) GetByID(ctx context.Context, id int64) (Article, error) {
	panic("unimplemented")
}

// GetByAuthor implements ArticleDAO.
func (m *MongoDBArticleDAO) GetByAuthor(
	ctx context.Context,
	uid int64,
	offset int,
	limit int,
) ([]Article, error) {
	panic("unimplemented")
}

// Insert implements ArticleDAO.
func (m *MongoDBArticleDAO) Insert(ctx context.Context, article Article) (int64, error) {
	article.ID = m.node.Generate().Int64()
	now := time.Now().UnixMilli()
	article.Ctime = now
	article.Utime = now
	_, err := m.coll.InsertOne(ctx, &article)
	return article.ID, err
}

// Sync implements ArticleDAO.
func (m *MongoDBArticleDAO) Sync(ctx context.Context, article Article) (int64, error) {
	var err error
	id := article.ID

	if id > 0 {
		err = m.UpdateByID(ctx, article)
	} else {
		id, err = m.Insert(ctx, article)
		article.ID = id
	}
	if err != nil {
		return 0, nil
	}

	// liveColl upsert
	err = m.Upsert(ctx, PublishedArticle(article))
	return id, err
}

// Transaction implements ArticleDAO.
func (m *MongoDBArticleDAO) Transaction(
	ctx context.Context,
	fn func(ctx context.Context, tx any) (any, error),
) (any, error) {
	sess, err := m.client.StartSession()
	if err != nil {
		return nil, err
	}
	defer func() {
		sess.EndSession(ctx)
	}()

	_, err = sess.WithTransaction(ctx, func(ctx context.Context) (interface{}, error) {
		return fn(ctx, m)
	})
	return nil, err
}

// UpdateByID implements ArticleDAO.
func (m *MongoDBArticleDAO) UpdateByID(ctx context.Context, article Article) error {
	now := time.Now().UnixMilli()
	filter := bson.M{
		"id":        article.ID,
		"author_id": article.AuthorID,
	}
	set := bson.M{
		"$set": bson.M{
			"title":   article.Title,
			"content": article.Content,
			"status":  article.Status,
			"utime":   now,
		},
	}
	res, err := m.coll.UpdateOne(ctx, filter, set)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return ErrArticleNotFound
	}
	return nil
}

// UpdateStatusByID implements ArticleDAO.
func (m *MongoDBArticleDAO) UpdateStatusByID(
	ctx context.Context,
	model any,
	userID int64,
	articleID int64,
	status uint8,
) error {
	panic("unimplemented")
}

// Upsert implements ArticleDAO.
func (m *MongoDBArticleDAO) Upsert(ctx context.Context, article PublishedArticle) error {
	now := time.Now().UnixMilli()
	article.Utime = now
	filter := bson.M{
		"id":        article.ID,
		"author_id": article.AuthorID,
	}
	set := bson.M{
		"$set": article,
		// set new attribut on onsert
		"$setOnInsert": bson.M{
			"ctime": now,
		},
	}

	res, err := m.liveColl.UpdateOne(ctx, filter, set, options.UpdateOne().SetUpsert(true))
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return ErrArticleNotFound
	}
	return nil
}

func NewMongoDBArticleDAO(client *mongo.Client, node *snowflake.Node) ArticleDAO {
	return &MongoDBArticleDAO{
		node:     node,
		client:   client,
		coll:     client.Database(DatabaseName).Collection(ArticleCollName),
		liveColl: client.Database(DatabaseName).Collection(PublishedArticleCollName),
	}
}

var _ ArticleDAO = &MongoDBArticleDAO{}
