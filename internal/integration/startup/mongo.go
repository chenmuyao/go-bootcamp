package startup

import (
	"context"
	"fmt"

	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
	"go.mongodb.org/mongo-driver/v2/event"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func InitMongoDB() *mongo.Client {
	monitor := &event.CommandMonitor{
		Started: func(ctx context.Context, cse *event.CommandStartedEvent) {
			fmt.Println(cse.Command)
		},
	}
	opts := options.Client().
		ApplyURI("mongodb://root:root@localhost:27017").
		SetMonitor(monitor)
	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}
	mdb := client.Database("wetravel")
	err = dao.InitCollection(mdb)
	if err != nil {
		panic(err)
	}
	return client
}
