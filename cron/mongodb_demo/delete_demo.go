package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type TimeBeforeCond struct {
	Before int64 `bson:"&lt"`
}
type DeleteCnd struct {
	beforeCond TimeBeforeCond `bson:"timeOption.startTime"`
}

func main() {
	var (
		ctx          context.Context
		client       *mongo.Client
		database     *mongo.Database
		collection   *mongo.Collection
		url          string
		err          error
		delCond      *DeleteCnd
		deleteResult *mongo.DeleteResult
	)
	url = "mongodb://127.0.0.1:27017"
	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	if client, err = mongo.Connect(ctx, options.Client().ApplyURI(url)); err != nil {
		fmt.Println(err)
		return
	}
	database = client.Database("cron")
	collection = database.Collection("log")
	delCond = &DeleteCnd{beforeCond: TimeBeforeCond{time.Now().Unix()}}

	if deleteResult, err = collection.DeleteMany(ctx, delCond); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(deleteResult.DeletedCount)
}
