package mongodblearn

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/event"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDBTestSuite struct {
	suite.Suite
	coll *mongo.Collection
}

func (s *MongoDBTestSuite) SetupSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	monitor := &event.CommandMonitor{
		Started: func(ctx context.Context, cse *event.CommandStartedEvent) {
			fmt.Println(cse.Command)
		},
	}
	opts := options.Client().
		ApplyURI("mongodb://root:root@localhost:27017").
		SetMonitor(monitor)
	client, err := mongo.Connect(opts)
	assert.NoError(s.T(), err)

	s.coll = client.Database("wetravel").Collection("articles")

	manyRes, err := s.coll.InsertMany(ctx, []any{
		Article{
			ID:       1,
			AuthorID: 1234,
		}, Article{
			ID:       2,
			AuthorID: 12345,
		},
	})
	assert.NoError(s.T(), err)
	s.T().Log("Insert", manyRes)
}

func (s *MongoDBTestSuite) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, _ = s.coll.DeleteMany(ctx, bson.D{})
	_ = s.coll.Indexes().DropAll(ctx)
}

func (s *MongoDBTestSuite) TestOr() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"$or": bson.A{
			bson.M{"id": 1},
			bson.M{"id": 2},
		},
	}

	// filter := bson.A{
	// 	bson.D{
	// 		bson.E{
	// 			Key:   "id",
	// 			Value: 123,
	// 		},
	// 	},
	// 	bson.D{
	// 		bson.E{
	// 			Key:   "id",
	// 			Value: 1234,
	// 		},
	// 	},
	// }

	res, err := s.coll.Find(ctx, filter)
	assert.NoError(s.T(), err)
	var articles []Article
	err = res.All(ctx, &articles)
	assert.NoError(s.T(), err)
	s.T().Log("Find or", articles)
}

func (s *MongoDBTestSuite) TestAnd() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"$and": bson.A{
			bson.M{"id": 1},
			bson.M{"author_id": 1234},
		},
	}

	res, err := s.coll.Find(ctx, filter)
	assert.NoError(s.T(), err)
	var articles []Article
	err = res.All(ctx, &articles)
	assert.NoError(s.T(), err)
	s.T().Log("Find and", articles)
}

func (s *MongoDBTestSuite) TestInWithProjection() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"id": bson.M{
			"$in": []int{1, 2, 3},
		},
	}

	res, err := s.coll.Find(ctx, filter, options.Find().SetProjection(bson.M{"id": 1}))
	assert.NoError(s.T(), err)
	var articles []Article
	err = res.All(ctx, &articles)
	assert.NoError(s.T(), err)
	s.T().Log("Find in", articles)
}

func (s *MongoDBTestSuite) TestIndexes() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	idxRes, err := s.coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"id": 1},
		Options: options.Index().SetUnique(true),
		// Options: options.Index().SetUnique(true).SetName("id_idx_1"),
	})
	assert.NoError(s.T(), err)
	s.T().Log("Index", idxRes)
}

func TestMongoDBQueries(t *testing.T) {
	suite.Run(t, &MongoDBTestSuite{})
}
