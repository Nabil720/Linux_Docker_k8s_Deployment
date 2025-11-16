package main

import (
	"log"
	"net/http"
	"employeeservice/database"
	"employeeservice/handlers"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func employeeHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	
	// Handle preflight OPTIONS request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Route based on the path
	switch r.URL.Path {
	case "/emp/employees":
		if r.Method == http.MethodGet {
			handlers.GetEmployees(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case "/emp/add-employee":
		if r.Method == http.MethodPost {
			handlers.AddEmployee(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case "/emp/delete-employee":
		if r.Method == http.MethodDelete {
			handlers.DeleteEmployee(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case "/emp/update-employee":
		if r.Method == http.MethodPut {
			handlers.UpdateEmployee(w, r)
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

	// Register the employee handler
	http.HandleFunc("/emp/", employeeHandler)
	
	// Root path handler for health check
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Employee Service is running"))
		}
	})

	log.Println("Employee Service running on port 5003")
	log.Fatal(http.ListenAndServe(":5003", nil))
}