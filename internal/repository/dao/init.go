package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"gorm.io/gorm"
)

func InitTable(db *gorm.DB) error {
	// NOTE: Not the best practice. Too risky. Strong dependency
	return db.AutoMigrate(&User{}, &SMSInfo{}, &Article{}, &PublishedArticle{})
}

func InitCollection(mdb *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	col := mdb.Collection(ArticleCollName)
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.M{"id": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{"author_id": 1},
		},
	})
	if err != nil {
		return err
	}

	liveColl := mdb.Collection(PublishedArticleCollName)
	_, err = liveColl.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.M{"id": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{"author_id": 1},
		},
	})
	if err != nil {
		return err
	}

	return nil
}
