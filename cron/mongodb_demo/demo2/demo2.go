package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type FindByJobName struct {
	JobName string `bson:"jobName"`
}

type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime   int64 `bson:"endTime"`
}

type LogRecord struct {
	JobName   string    `bson:"jobName"`
	Command   string    `bson:"command"`
	Err       string    `bson:"err"`
	Content   string    `bson:"content"`
	TimePoint TimePoint `bson:"timePoint"`
}

func main() {
	var (
		URl        string
		client     *mongo.Client
		err        error
		database   *mongo.Database
		ctx        context.Context
		collection *mongo.Collection
		cond       *FindByJobName
		result     *mongo.Cursor
		findOpt    *options.FindOptions
		record     *LogRecord
	)
	URl = "mongodb://127.0.0.1:27017"
	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	if client, err = mongo.Connect(ctx, options.Client().ApplyURI(URl)); err != nil {
		fmt.Println(err)
		return
	}
	findOpt =&options.FindOptions{}
	database = client.Database("cron")
	collection = database.Collection("log")
	cond = &FindByJobName{
		JobName: "job10",
	}
	if result,err = collection.Find(ctx,cond,findOpt.SetSkip(0),findOpt.SetLimit(2));err!=nil {
		fmt.Println(err)
		return
	}
	//释放资源
	defer result.Close(context.TODO())

	for result.Next(context.TODO()) {
		record =  &LogRecord{}
		if err = result.Decode(record); err!=nil {
			fmt.Println(err)
			return
		}
		fmt.Println(*record)
	}


}
