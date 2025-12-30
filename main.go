package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient       *mongo.Client
	moviesCollections *mongo.Collection
)

type Movies struct {
	ID     primitive.ObjectID `json:"id" bson:"id"`
	Title  string             `json:"title" bson:"title"`
	Rating float64            `json:"rating" bson:"rating"`
}

func main() {
	initMongoDB()
	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Get("/getMovies", handleAllMovies)
		r.Post("/createMovie", handleCreateMovie)
		r.Put("/updateMovie", handleUpdateMovie)
	})
	http.ListenAndServe(":4839", r)
	fmt.Print("server running on port 4839")
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

}

func handleCreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}

func handleUpdateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}
