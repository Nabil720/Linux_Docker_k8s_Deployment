package main

import (
	"log"
	"net/http"
	"studentservice/database" 
	"studentservice/handlers"

	"go.elastic.co/apm/v2"
)

// Vault থেকে Environment Variables লোড করার function
func loadEnvFromVault() {
	// Vault থেকে secrets /vault/secrets/config ফাইলে থাকে
	// তারা automagically environment variables হয়ে যায়
	log.Println("Vault secrets loaded automatically via sidecar")
}

func initAPM() {
	loadEnvFromVault()
	
	// APM initialization
	if apm.DefaultTracer().Active() {
		log.Println("APM initialized for Student Service")
	} else {
		log.Println("APM not active - using environment variables")
	}
}

// SIMPLIFIED APM middleware
func apmMiddleware(handler http.HandlerFunc, operationName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tracer := apm.DefaultTracer()
		if tracer == nil || !tracer.Active() {
			handler(w, r)
			return
		}
		
		tx := tracer.StartTransaction(operationName, "request")
		defer tx.End()
		
		ctx := apm.ContextWithTransaction(r.Context(), tx)
		req := r.WithContext(ctx)
		handler(w, req)
	}
}


func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func main() {
	// Initialize APM
	initAPM()

	// Database connection
	if err := database.Connect(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Student Routes
	http.HandleFunc("/std/add-student", func(w http.ResponseWriter, r *http.Request) {   // Here using this  "/std/add-student" path , we POST in  the Add-student service 
		enableCors(w)
		if r.Method == http.MethodOptions {
			return
		}
		apmMiddleware(handlers.AddStudent, "POST /add-student")(w, r)
	})

	http.HandleFunc("/std/students", func(w http.ResponseWriter, r *http.Request) {         // Here using this  "/std/students" path , we GET students  
		enableCors(w)
		enableCors(w)
		if r.Method == http.MethodOptions {
			return
		}
		apmMiddleware(handlers.GetStudents, "GET /students")(w, r)
	})

	http.HandleFunc("/std/delete-student", func(w http.ResponseWriter, r *http.Request) {    // Here using this  "/std/delete-student" path , we DELETE in  the Students  service 
		enableCors(w)
		enableCors(w)
		if r.Method == http.MethodOptions {
			return
		}
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		apmMiddleware(handlers.DeleteStudent, "DELETE /delete-student")(w, r)
	})

	http.HandleFunc("/std/update-student", func(w http.ResponseWriter, r *http.Request) {  // Here using this  "/std/update-student" path , we PUT in  the Students service 
		enableCors(w)
		if r.Method == http.MethodOptions {
			return
		}
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		apmMiddleware(handlers.UpdateStudent, "PUT /update-student")(w, r)
	})

	log.Println("Student Service running on port 5001")
	log.Fatal(http.ListenAndServe(":5001", nil))
}