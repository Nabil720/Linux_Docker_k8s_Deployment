package main

import (
	"log"
	"net/http"
	"teacherservice/database" 
	"teacherservice/handlers"
)

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func main() {
	// Database connection
	if err := database.Connect(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Teacher Routes
	http.HandleFunc("/add-teacher", func(w http.ResponseWriter, r *http.Request) {
		enableCors(w)
		if r.Method == http.MethodOptions {
			return
		}
		handlers.AddTeacher(w, r)
	})

	http.HandleFunc("/teachers", func(w http.ResponseWriter, r *http.Request) {
		enableCors(w)
		if r.Method == http.MethodOptions {
			return
		}
		handlers.GetTeachers(w, r)
	})

	http.HandleFunc("/delete-teacher", func(w http.ResponseWriter, r *http.Request) {
		enableCors(w)
		if r.Method == http.MethodOptions {
			return
		}
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handlers.DeleteTeacher(w, r)
	})

	http.HandleFunc("/update-teacher", func(w http.ResponseWriter, r *http.Request) {
		enableCors(w)
		if r.Method == http.MethodOptions {
			return
		}
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handlers.UpdateTeacher(w, r)
	})

	log.Println("Teacher Service running on port 5002")
	log.Fatal(http.ListenAndServe(":5002", nil))
}