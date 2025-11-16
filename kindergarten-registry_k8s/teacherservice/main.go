package main

import (
	"log"
	"net/http"
	"teacherservice/database"
	"teacherservice/handlers"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func teacherHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	
	// Handle preflight OPTIONS request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Route based on the path
	switch r.URL.Path {
	case "/tech/teachers":
		if r.Method == http.MethodGet {
			handlers.GetTeachers(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case "/tech/add-teacher":
		if r.Method == http.MethodPost {
			handlers.AddTeacher(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case "/tech/delete-teacher":
		if r.Method == http.MethodDelete {
			handlers.DeleteTeacher(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case "/tech/update-teacher":
		if r.Method == http.MethodPut {
			handlers.UpdateTeacher(w, r)
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

	// Register the teacher handler
	http.HandleFunc("/tech/", teacherHandler)
	
	// Root path handler for health check
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Teacher Service is running"))
		}
	})

	log.Println("Teacher Service running on port 5002")
	log.Fatal(http.ListenAndServe(":5002", nil))
}