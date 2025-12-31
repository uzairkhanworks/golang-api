package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient       *mongo.Client
	moviesCollections *mongo.Collection
)

type Movies struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title  string             `json:"title" bson:"title"`
	Rating float64            `json:"rating" bson:"rating"`
}

func main() {
	initMongoDB()
	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Get("/getMovies", handleAllMovies)
		r.Post("/createMovie", handleCreateMovie)
		r.Put("/updateMovie/{id}", handleUpdateMovie)
		r.Delete("/deleteMovie/{id}", handleDeleteMovie)
	})
	fmt.Print("server running on port 4839")
	http.ListenAndServe(":4839", r)
}

func initMongoDB() {
	//need to write mongo db connect logic
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("could not connect to mongoDB")
	}
	mongoClient = client
	moviesCollections = client.Database("movies").Collection("moviesCollection")
	fmt.Print("Successfully connected to DB")

}
func handleAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var movies []Movies
	cursor, err := moviesCollections.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &movies)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if movies == nil {
		movies = []Movies{}
	}
	json.NewEncoder(w).Encode(movies)
}

func handleCreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var payloadMovie Movies
	if err := json.NewDecoder(r.Body).Decode(&payloadMovie); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if payloadMovie.Title == "" {
		http.Error(w, "New Movie must have a title", http.StatusBadRequest)
	}
	result, err := moviesCollections.InsertOne(ctx, payloadMovie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var insertedMovie Movies
	err = moviesCollections.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&insertedMovie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(insertedMovie)
}

func handleUpdateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	id := chi.URLParam(r, "id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}
	var updateMovie Movies
	err = json.NewDecoder(r.Body).Decode(&updateMovie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	update := bson.M{
		"$set": bson.M{
			"title":  updateMovie.Title,
			"rating": updateMovie.Rating,
		},
	}
	filter := bson.M{"_id": objectID}
	result := moviesCollections.FindOneAndUpdate(ctx, filter, update)
	var updatedMovieResult Movies
	if err := result.Decode(&updatedMovieResult); err != nil {
		http.Error(w, "Movie not found or update failed", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(updatedMovieResult)

}

func handleDeleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	id := chi.URLParam(r, "id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "invalid id provided", http.StatusBadRequest)
	}
	result := moviesCollections.FindOneAndDelete(ctx, bson.M{"_id": objectID})
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Successfully removed the user",
		"result":  result,
	})
}
