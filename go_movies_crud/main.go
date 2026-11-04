package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
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

var (
	mu     sync.RWMutex
	movies []Movie // dummy DB
)

func getMovies(w http.ResponseWriter, r *http.Request) {
	var responseMovies []Movie
	mu.RLock()
	responseMovies = append([]Movie(nil), movies...)
	mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseMovies); err != nil {
		log.Printf("encode movies response: %v", err)
		return
	}

}
func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	found := false
	var responseMovie Movie
	mu.RLock()
	for _, movie := range movies {
		if movie.Id == params["id"] {
			found = true
			responseMovie = movie
		}
	}
	mu.RUnlock()

	if found {
		if err := json.NewEncoder(w).Encode(responseMovie); err != nil {
			log.Printf("encode movie response: %v", err)
			return
		}
		return
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

	var responseMovies []Movie
	mu.Lock()
	movies = append(movies, movie)
	responseMovies = append([]Movie(nil), movies...)
	mu.Unlock()

	//new movies
	if err := json.NewEncoder(w).Encode(responseMovies); err != nil {
		log.Printf("encode movies response: %v", err)
		return
	}
}

func updateMovie(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var new_movie Movie
	err := json.NewDecoder(r.Body).Decode(&new_movie)
	if err != nil {
		http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
		return
	}

	mu.Lock()
	found := false
	var responseMovies []Movie

	for i, movie := range movies {
		if movie.Id == params["id"] {
			movies = append(movies[:i], movies[i+1:]...)
			movies = append(movies, new_movie)
			responseMovies = append([]Movie(nil), movies...)
			found = true
			break
		}
	}
	mu.Unlock()

	if !found {
		http.Error(w, "Movie doesn't exist", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(responseMovies); err != nil {
		log.Printf("encode movies response: %v", err)
		return
	}
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	found := false
	var responseMovies []Movie
	mu.Lock()
	for i, movie := range movies {
		if movie.Id == params["id"] {
			movies = append(movies[:i], movies[i+1:]...)
			responseMovies = append([]Movie(nil), movies...)
			found = true
			break
		}
	}
	mu.Unlock()

	if found {
		if err := json.NewEncoder(w).Encode(responseMovies); err != nil {
			log.Printf("encode movies response: %v", err)
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

	fmt.Println("Starting Server at PORT 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
