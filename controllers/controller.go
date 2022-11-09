package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"crud/model"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongodb helpers
const connectionString = "mongodb+srv://prajjawal0311:prajjawal0311@cluster0.qncwagp.mongodb.net/?retryWrites=true&w=majority"
const dbName = "netflix"
const colName = "watchlist"

var collection *mongo.Collection

// special method : will only runs once at the start of the execution
func init() {
	//client options
	clientOptions := options.Client().ApplyURI(connectionString)
	//connect to mongodb
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("mongodb connection success")
	collection = client.Database(dbName).Collection(colName)
	fmt.Println("collection reference is ready!!")
}

// insert 1 record
func insertOneMovie(movie model.Netflix) {
	ins, err := collection.InsertOne(context.Background(), movie)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("inserted 1 movie in db with id:", ins.InsertedID)
}

// update 1 record
func updateOneMovie(movieid string) {
	id, _ := primitive.ObjectIDFromHex(movieid)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"watched": true}}
	res, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("modified count:", res.ModifiedCount)
}

// delete 1 record
func deleteOneMovie(movieid string) {
	id, _ := primitive.ObjectIDFromHex(movieid)
	filter := bson.M{"_id": id}
	res, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("delete count :", res.DeletedCount)
}

// fetch all data
func getAllMovies() []primitive.M {
	cursor, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	var movies []primitive.M
	for cursor.Next(context.Background()) {
		var movie bson.M
		err = cursor.Decode(&movie)
		if err != nil {
			log.Fatal(err)
		}
		movies = append(movies, movie)
	}
	defer cursor.Close(context.Background())
	return movies
}

//---------------------------------controllers----------------------------------------

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var movie model.Netflix
	json.NewDecoder(r.Body).Decode(&movie)
	insertOneMovie(movie)
	json.NewEncoder(w).Encode(movie)
}
func GetMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	allMovies := getAllMovies()
	json.NewEncoder(w).Encode(allMovies)
}
func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	params := mux.Vars(r)
	updateOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}
func DeleteMovie(w http.ResponseWriter, r *http.Request){
	w.Header().Set("content-type", "application/json")
	params := mux.Vars(r)
	deleteOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}