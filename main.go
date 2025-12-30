package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
}
func handleAllMovies(w http.ResponseWriter, r *http.Request) {

}

func handleCreateMovie(w http.ResponseWriter, r *http.Request) {

}

func handleUpdateMovie(w http.ResponseWriter, r *http.Request) {

}
