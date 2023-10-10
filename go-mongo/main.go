package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type Movie struct {
    Name string `bson:"name"`
    Year uint16 `bson:"year"`
    Directors []string `bson:"directors"`
    Writers []string `bson:"writers"`
    BoxOffice `bson:"boxOffice"`
}


type BoxOffice struct{
    Budget uint64 `bson:"budget"`
    Gross uint64 `bson:"gross"` 
}

func main()  {
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(context.TODO(), clientOptions)

    if err != nil {
        panic(err)
    }
    err = client.Ping(context.TODO(), nil)

    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Connected to mongodb success !!")
    collection := client.Database("appDB").Collection("movies")
    darkKnight := Movie{
        Name: "The Dark Knight",
        Year: 2023,
        Directors: []string{"Christopher Nolan"},
        Writers: []string{"Jonathan Nolan", "Christopher Nolan"},
        BoxOffice: BoxOffice{
            Budget: 185000000,
            Gross: 533316061,
        },
    }
    _, err = collection.InsertOne(context.TODO(), darkKnight) 
    if err != nil {
        log.Fatal(err)
    }
    queryResult := &Movie{}
    filter := bson.M{"boxOffice.budget": bson.M{"$gt": 150000000}}
    result := collection.FindOne(context.TODO(), filter)
    err = result.Decode(queryResult)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Movie ==== ",queryResult)
    err = client.Disconnect(context.TODO())
    if err != nil {
        panic(err)
    }
    fmt.Println("Dissconnected from MONGO777")
}

