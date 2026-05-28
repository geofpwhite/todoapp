package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"todo/gom"
	"todo/records"
	// "github.com/gorilla/mux" // Import the mux package
)

func main() {
	mux := http.NewServeMux()

	rh := records.NewRecordHandler()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		gom.Home(w, rh)
	})
	mux.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Unable to read request body", http.StatusBadRequest)
			}
			task := string(body)
			rh.AddRecord(task)
		case "OPTIONS":
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		}
	})
	mux.HandleFunc("/complete/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		num, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		rh.MarkRecordAsCompleted(num)
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/completed", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/?tab=completed", http.StatusFound)
	})

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Starting server at :%s\n", port)

	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
}
