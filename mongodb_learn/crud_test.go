package mongodblearn

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/event"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Article struct {
	ID      int64  `bson:"id,omitempty"`
	Title   string `bson:"title,omitempty"`
	Content string `bson:"content,omitempty"`

	AuthorID int64 `bson:"author_id,omitempty"`
	Status   uint8 `bson:"status,omitempty"`
	Ctime    int64 `bson:"ctime,omitempty"`
	Utime    int64 `bson:"utime,omitempty"`
}

func TestMondoDB(t *testing.T) {
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
	assert.NoError(t, err)

	coll := client.Database("wetravel").Collection("articles")

	// C
	insertRes, err := coll.InsertOne(ctx, Article{
		ID:       1,
		Title:    "my title",
		Content:  "my content",
		AuthorID: 123,
	})
	assert.NoError(t, err)
	oid := insertRes.InsertedID.(bson.ObjectID)
	t.Log("insert ID ", oid)

	// R
	// filter := bson.D{bson.E{Key: "id", Value: 1}}
	filter := bson.M{
		"id": 1,
	}
	findRes := coll.FindOne(ctx, filter)
	if findRes.Err() == mongo.ErrNoDocuments {
		t.Log("not found")
	} else {
		var article Article
		err = findRes.Decode(&article)
		assert.NoError(t, err)
		t.Log(article)
	}

	// U
	updateFilter := bson.D{bson.E{Key: "id", Value: 1}}
	set := bson.D{bson.E{
		Key: "$set",
		Value: bson.E{
			Key:   "title",
			Value: "New title",
		},
	}}
	updateOnRes, err := coll.UpdateOne(ctx, updateFilter, set)
	assert.NoError(t, err)
	t.Log("update document amount", updateOnRes.ModifiedCount)

	updateManyRes, err := coll.UpdateMany(ctx, updateFilter, bson.D{bson.E{
		Key:   "$set",
		Value: Article{Content: "New content"},
	}})
	assert.NoError(t, err)
	t.Log("update content amount", updateManyRes.ModifiedCount)

	// D
	deleteFilter := bson.D{bson.E{Key: "id", Value: 1}}
	deleteManyRes, err := coll.DeleteMany(ctx, deleteFilter)
	assert.NoError(t, err)
	t.Log("delete content amount", deleteManyRes.DeletedCount)
}
