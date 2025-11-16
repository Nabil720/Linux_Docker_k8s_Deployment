package main

import (
	"log"
	"net/http"
	"employeeservice/database" 
	"employeeservice/handlers"
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

	// Employee Routes
	http.HandleFunc("/emp/add-employee", func(w http.ResponseWriter, r *http.Request) {   // Here using this  "/emp/add-employee" path , we POST in  the Add-employee service 
		enableCors(w)
		if r.Method == http.MethodOptions {
			return
		}
		handlers.AddEmployee(w, r)
	})

	http.HandleFunc("/emp/employees", func(w http.ResponseWriter, r *http.Request) {    // Here using this  "/emp/employees" path , we GET employees  
		enableCors(w)
		if r.Method == http.MethodOptions {
			return
		}
		handlers.GetEmployees(w, r)
	})

	http.HandleFunc("/emp/delete-employee", func(w http.ResponseWriter, r *http.Request) {   // Here using this  "/emp/delete-employee" path , we DELETE in  the Employees service 
		enableCors(w)
		enableCors(w)
		if r.Method == http.MethodOptions {
			return
		}
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handlers.DeleteEmployee(w, r)
	})

	http.HandleFunc("/emp/update-employee", func(w http.ResponseWriter, r *http.Request) {    // Here using this  "/emp/update-employee" path , we PUT in  the Employees service 
		enableCors(w)
		if r.Method == http.MethodOptions {
			return
		}
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handlers.UpdateEmployee(w, r)
	})

	log.Println("Employee Service running on port 5003")
	log.Fatal(http.ListenAndServe(":5003", nil))
}