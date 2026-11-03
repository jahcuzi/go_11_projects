package main

import (
	"encoding/json"
	"log"
	"net/http"
	"uuid"

	"github.com/gorilla/mux"
)

type Movie struct {
	Id       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie // dummy DB

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(movies); err != nil {
		log.Printf("encode movies response: %v", err)
		return
	}

}
func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, movie := range movies {
		if movie.Id == params["id"] {
			if err := json.NewEncoder(w).Encode(movie); err != nil {
				log.Printf("encode movie response: %v", err)
				return
			}
			return
		}
	}

	http.Error(w, "Movie not found", http.StatusNotFound)
}
func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
		return
	}
	movie.Id = uuid.NewV7().String()
	movies = append(movies, movie)
	//new movies
	if err := json.NewEncoder(w).Encode(movies); err != nil {
		log.Printf("encode movies response: %v", err)
		return
	}
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var movie Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
		return
	}

	for i, m := range movies {
		if m.Id == params["id"] {
			movies = append(movies[:i], movies[i+1:]...)
			movie.Id = params["id"]
			movies = append(movies, movie)
			if err := json.NewEncoder(w).Encode(movies); err != nil {
				log.Printf("encode movies response: %v", err)
				return
			}
			return
		}
	}

	http.Error(w, "Movie not found", http.StatusNotFound)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for i, movie := range movies {
		if movie.Id == params["id"] {
			movies = append(movies[:i], movies[i+1:]...)
			if err := json.NewEncoder(w).Encode(movies); err != nil {
				log.Printf("encode movies response: %v", err)
				return
			}
			return
		}
	}

	http.Error(w, "Movie doesn't exist", http.StatusNotFound)

}

func main() {
	r := mux.NewRouter()
	movies = append(movies, Movie{Id: "1", Isbn: "121851", Title: "Movie 1", Director: &Director{Firstname: "f_name_1", Lastname: "l_name_1"}})
	movies = append(movies, Movie{Id: "2", Isbn: "15454851", Title: "Movie 2", Director: &Director{Firstname: "f_name_2", Lastname: "l_name_2"}})
	r.HandleFunc("/movies", getMovies).Methods(http.MethodGet)
	r.HandleFunc("/movie/{id}", getMovie).Methods(http.MethodGet)
	r.HandleFunc("/movie", createMovie).Methods(http.MethodPost)
	r.HandleFunc("/movie/{id}", updateMovie).Methods(http.MethodPut)
	r.HandleFunc("/movie/{id}", deleteMovie).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(":8080", r))
}
