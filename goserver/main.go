package main

import (
	"fmt"
	"log"
	"net/http"
)

func handleHome(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/home" {
		http.Error(w, "404 NOt found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "method not supported", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "Hiii, AMan\n")
}

func handleForm(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not supported", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("ParseForm() error: %v", err), http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	address := r.FormValue("address")
	fmt.Fprintln(w, "POST request successful")
	fmt.Fprintln(w, "name is", name)
	fmt.Fprintln(w, "Address is", address)

}

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/home", handleHome)
	http.HandleFunc("/form", handleForm)

	fmt.Println("Starting Server at PORT :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
