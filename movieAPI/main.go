package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type DB struct{
    collection *mongo.Collection
}

type Movie struct {
    ID interface{} `json:"id" bson: "_id,omitempty"` 
    Name string `json:"name" bson:"name"`
    Year uint16 `json:"year" bson:"year"`
    Directors []string `json:"directors" bson:"directors"`
    Writers []string `json:"writers" bson:"writers"`
    BoxOffice BoxOffice `json:"boxOffice" bson:"boxOffice"`
}

type BoxOffice struct{
    Budget uint64 `json:"budget" bson:"budget"`
    Gross uint64 `json:"gross" bson:"gross"`
}


func (db *DB) GetMovie(w http.ResponseWriter, r *http.Request)  {
    vars := mux.Vars(r)
    var movie Movie
    objectID, err := primitive.ObjectIDFromHex(vars["id"])
    filter := bson.M{"_id": objectID}
    err = db.collection.FindOne(context.TODO(), filter).Decode(&movie)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
    }else {
        w.Header().Set("Content-Type", "application/json")
        response, _ := json.Marshal(movie)
        w.WriteHeader(http.StatusOK)
        w.Write(response)
    }
    
}
func (db *DB) AddMovie(w http.ResponseWriter, r *http.Request)  {
    var movie Movie
    postBody, _ := io.ReadAll(r.Body)
    json.Unmarshal(postBody, &movie)
    result, err := db.collection.InsertOne(context.TODO(), movie)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
    }else {
        w.Header().Set("Content-Type", "application/json")
        response, _ := json.Marshal(result)
        w.WriteHeader(http.StatusOK)
        w.Write(response)
    }
    
}

func (db *DB) removeMovie(w http.ResponseWriter, r *http.Request)  {
    vars := mux.Vars(r)
    objectID, err := primitive.ObjectIDFromHex(vars["id"])
    filter := bson.M{"_id": objectID}
    _ ,err = db.collection.DeleteOne(context.TODO(), filter)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
    }else {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Deleted"))
    }
}

func main()  {
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        panic(err)
    }
    defer client.Disconnect(context.TODO())

    collection := client.Database("appDB").Collection("movies")
    db := &DB{collection: collection}

    r := mux.NewRouter()
    r.HandleFunc("/v1/movies/{id:[a-zA-Z0-9]*}", db.GetMovie).Methods("GET")
    r.HandleFunc("/v1/movies", db.AddMovie).Methods("POST")
    r.HandleFunc("/v1/movies/{id:[a-zA-Z0-9]*", db.removeMovie).Methods("DELETE")
    srv := &http.Server{
        Handler: r,
        Addr: "127.0.0.1:8000",
        WriteTimeout: 15 * time.Second,
        ReadTimeout: 15 * time.Second,
    }
    log.Fatal(srv.ListenAndServe())
}
