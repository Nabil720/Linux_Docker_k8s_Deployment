package main

import (
	"log"
	"net/http"
	"studentservice/database"
	"studentservice/handlers"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func studentHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	
	// Handle preflight OPTIONS request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Route based on the path
	switch r.URL.Path {
	case "/std/students":
		if r.Method == http.MethodGet {
			handlers.GetStudents(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case "/std/add-student":
		if r.Method == http.MethodPost {
			handlers.AddStudent(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case "/std/delete-student":
		if r.Method == http.MethodDelete {
			handlers.DeleteStudent(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case "/std/update-student":
		if r.Method == http.MethodPut {
			handlers.UpdateStudent(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	default:
		http.NotFound(w, r)
	}
}

func main() {
	// Database connection
	if err := database.Connect(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Register the student handler
	http.HandleFunc("/std/", studentHandler)
	
	// Root path handler for health check
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Student Service is running"))
		}
	})

	log.Println("Student Service running on port 5001")
	log.Fatal(http.ListenAndServe(":5001", nil))
}