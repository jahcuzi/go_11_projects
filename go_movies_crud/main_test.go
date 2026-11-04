package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/gorilla/mux"
)

func testServer() *mux.Router {
	mr := mux.NewRouter()
	mr.HandleFunc("/movies", getMovies).Methods(http.MethodGet)
	mr.HandleFunc("/movie/{id}", getMovie).Methods(http.MethodGet)
	mr.HandleFunc("/movie", createMovie).Methods(http.MethodPost)
	mr.HandleFunc("/movie/{id}", updateMovie).Methods(http.MethodPut)
	mr.HandleFunc("/movie/{id}", deleteMovie).Methods(http.MethodDelete)
	return mr
}

func setMovies(t *testing.T, testMovies []Movie) {
	t.Helper()

	mu.Lock()
	oldMovies := append([]Movie(nil), movies...)
	movies = append([]Movie(nil), testMovies...)
	mu.Unlock()

	t.Cleanup(func() {
		mu.Lock()
		movies = oldMovies
		mu.Unlock()
	})
}

func TestMovieNotFound(t *testing.T) {
	setMovies(t, []Movie{{Id: "1", Title: "meh"}})

	req := httptest.NewRequest(http.MethodGet, "/movie/5", nil)
	rec := httptest.NewRecorder()

	router := testServer()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("status =%d, want %d", rec.Code, http.StatusNotFound)
	}

}
func TestGetMovies(t *testing.T) {
	setMovies(t, []Movie{
		{Id: "1", Title: "First"},
		{Id: "2", Title: "Second"},
	})

	req := httptest.NewRequest(http.MethodGet, "/movies", nil)
	rec := httptest.NewRecorder()

	testServer().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body: %s",
			rec.Code, http.StatusOK, rec.Body.String())
	}

	var got []Movie
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	want := []Movie{
		{Id: "1", Title: "First"},
		{Id: "2", Title: "Second"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v movies, want %v", got, want)
	}
}

func TestGetMovie(t *testing.T) {
	setMovies(t, []Movie{
		{Id: "1", Title: "First"},
		{Id: "2", Title: "Second"},
	})

	req := httptest.NewRequest(http.MethodGet, "/movie/1", nil)
	rec := httptest.NewRecorder()

	testServer().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Status code not OK, got %d, want, %d; body:%s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var got Movie
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	want := Movie{Id: "1", Title: "First"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("movie = %+v, want %+v", got, want)
	}

}

func TestCreateMovie(t *testing.T) {
	setMovies(t, nil)

	body := `{"isbn":"12334","title":"movie 1","director":{"firstname":"fd","lastname":"ld"}}`
	req := httptest.NewRequest(http.MethodPost, "/movie", strings.NewReader(body))
	rec := httptest.NewRecorder()

	testServer().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Status: %d, want %d; body %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var responseMovies []Movie
	if err := json.NewDecoder(rec.Body).Decode(&responseMovies); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	mu.RLock()
	gotMovies := append([]Movie(nil), movies...)
	mu.RUnlock()

	if len(gotMovies) != 1 {
		t.Fatalf("Got movies length=%d, want 1", len(gotMovies))
	}

	got := gotMovies[0]
	if got.Id == "" {
		t.Fatal("created movie has an empty id")
	}
	want := Movie{
		Id:       got.Id,
		Isbn:     "12334",
		Title:    "movie 1",
		Director: &Director{Firstname: "fd", Lastname: "ld"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("stored movie = %+v, want %+v", got, want)
	}
	if !reflect.DeepEqual(responseMovies, []Movie{want}) {
		t.Fatalf("response movies = %+v, want %+v", responseMovies, []Movie{want})
	}
}

func TestUpdateMovie(t *testing.T) {
	setMovies(t, []Movie{{
		Id:    "1",
		Title: "title 1",
		Isbn:  "isbn 1",
		Director: &Director{
			Firstname: "fname 1",
			Lastname:  "lname 1",
		},
	}})

	body := `{"id": "1", "title":"title 2","isbn":"isbn 2", "director":{"firstname":"fname 2", "lastname":"lname 2"}}`
	req := httptest.NewRequest(http.MethodPut, "/movie/1", strings.NewReader(body))
	rec := httptest.NewRecorder()
	testServer().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Status: %d, want: %d; body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	//decode resp
	var responseMovies []Movie
	if err := json.NewDecoder(rec.Body).Decode(&responseMovies); err != nil {
		t.Fatalf("Decode error: %v", err)
	}

	//get got from Dummy DB
	mu.RLock()
	gotMovies := append([]Movie(nil), movies...)
	mu.RUnlock()

	//check in got movies for update
	want := Movie{
		Id:    "1",
		Title: "title 2",
		Isbn:  "isbn 2",
		Director: &Director{
			Firstname: "fname 2",
			Lastname:  "lname 2",
		},
	}

	if len(gotMovies) != 1 {
		t.Fatalf("Wrong length, got %d, want 1", len(gotMovies))
	}

	got := gotMovies[0]

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Wrong dummy DB, want %v, got: %v, body: %s", want, got, rec.Body.String())
	}

	//check in response movies

	if len(responseMovies) != 1 {
		t.Fatalf("Wrong length, got %d, want 1", len(responseMovies))
	}

	got = responseMovies[0]

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Wrong Response Movies, want %v, got: %v, body: %s", want, got, rec.Body.String())
	}

}

func TestDeleteMovie(t *testing.T) {
	setMovies(t, []Movie{{
		Id:    "1",
		Title: "title",
		Isbn:  "isbn",
		Director: &Director{
			Firstname: "fname",
			Lastname:  "lname",
		},
	}})

	req := httptest.NewRequest(http.MethodDelete, "/movie/1", nil)
	rec := httptest.NewRecorder()

	testServer().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Server Code Error: got %d, want %d; body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var responseMovies []Movie
	if err := json.NewDecoder(rec.Body).Decode(&responseMovies); err != nil {
		t.Fatalf("decode error %v", err)
	}

	mu.RLock()
	gotMovies := append([]Movie(nil), movies...)
	mu.RUnlock()

	toBeDeleted := Movie{
		Id:    "1",
		Title: "title",
		Isbn:  "isbn",
		Director: &Director{
			Firstname: "fname",
			Lastname:  "lname",
		},
	}

	//check response Movies
	for _, movie := range responseMovies {
		if movie.Id == toBeDeleted.Id {
			t.Fatalf("Movie of id: %s is presend in response", toBeDeleted.Id)
		}
	}

	//check got movies
	for _, movie := range gotMovies {
		if movie.Id == toBeDeleted.Id {
			t.Fatalf("MOvie of id: %s is present in DB", toBeDeleted.Id)
		}
	}

}

func TestCreateMovieInvalidJSON(t *testing.T) {
	initialMovies := []Movie{{Id: "1", Title: "Existing movie"}}
	setMovies(t, initialMovies)

	req := httptest.NewRequest(http.MethodPost, "/movie", strings.NewReader(`{"title":`))
	rec := httptest.NewRecorder()
	testServer().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body: %s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}

	mu.RLock()
	gotMovies := append([]Movie(nil), movies...)
	mu.RUnlock()
	if !reflect.DeepEqual(gotMovies, initialMovies) {
		t.Errorf("movies after invalid create = %+v, want unchanged %+v", gotMovies, initialMovies)
	}
}

func TestUpdateMovieInvalidJSON(t *testing.T) {
	initialMovies := []Movie{{Id: "1", Title: "Existing movie"}}
	setMovies(t, initialMovies)

	req := httptest.NewRequest(http.MethodPut, "/movie/1", strings.NewReader(`{"title":`))
	rec := httptest.NewRecorder()
	testServer().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body: %s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}

	mu.RLock()
	gotMovies := append([]Movie(nil), movies...)
	mu.RUnlock()
	if !reflect.DeepEqual(gotMovies, initialMovies) {
		t.Errorf("movies after invalid update = %+v, want unchanged %+v", gotMovies, initialMovies)
	}
}

func TestUpdateMovieNotFound(t *testing.T) {
	initialMovies := []Movie{{Id: "1", Title: "Existing movie"}}
	setMovies(t, initialMovies)

	body := `{"id":"2","title":"Updated movie"}`
	req := httptest.NewRequest(http.MethodPut, "/movie/2", strings.NewReader(body))
	rec := httptest.NewRecorder()
	testServer().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d; body: %s", rec.Code, http.StatusNotFound, rec.Body.String())
	}

	mu.RLock()
	gotMovies := append([]Movie(nil), movies...)
	mu.RUnlock()
	if !reflect.DeepEqual(gotMovies, initialMovies) {
		t.Errorf("movies after update for missing ID = %+v, want unchanged %+v", gotMovies, initialMovies)
	}
}

func TestDeleteMovieNotFound(t *testing.T) {
	initialMovies := []Movie{{Id: "1", Title: "Existing movie"}}
	setMovies(t, initialMovies)

	req := httptest.NewRequest(http.MethodDelete, "/movie/2", nil)
	rec := httptest.NewRecorder()
	testServer().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d; body: %s", rec.Code, http.StatusNotFound, rec.Body.String())
	}

	mu.RLock()
	gotMovies := append([]Movie(nil), movies...)
	mu.RUnlock()
	if !reflect.DeepEqual(gotMovies, initialMovies) {
		t.Errorf("movies after delete for missing ID = %+v, want unchanged %+v", gotMovies, initialMovies)
	}
}

func TestConcurrentCreateMovies(t *testing.T) {
	setMovies(t, nil)
	const REQUEST_COUNT = 1000
	router := testServer()
	var wg sync.WaitGroup
	errs := make(chan string, REQUEST_COUNT)

	for range REQUEST_COUNT {

		wg.Go(func() {

			body := `{"isbn":"000", "title":"concurrent test", "director":{"firstname":"fd", "lastname":"ld"}}`
			req := httptest.NewRequest(http.MethodPost, "/movie", strings.NewReader(body))
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				errs <- fmt.Sprintf("status =%d, want %d", rec.Code, http.StatusOK)
			}
		})
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Error(err)
	}

	mu.RLock()
	got := len(movies)
	mu.RUnlock()

	if REQUEST_COUNT != got {
		t.Errorf("stored %d movies, want %d", got, REQUEST_COUNT)
	}
}
