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
	os.Chdir("C:/Users/geoffrey/projects/todoapp/backend")
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
		gom.Completed(w, rh)
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("backend/static"))))

	// Start the server
	fmt.Println("Starting server at :8080")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
}
