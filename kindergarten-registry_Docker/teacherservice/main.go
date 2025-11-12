package main

import (
	"log"
	"net/http"
	"teacherservice/database" 
	"teacherservice/handlers"
	
	_ "go.elastic.co/apm/v2"
)

// APM configuration
func initAPM() {
	// APM server configuration
	err := apm.InitDefaultTracer(
		apm.WithServiceName("teacher-service"),
		apm.WithServiceVersion("1.0.0"),
		apm.WithServiceEnvironment("development"),
		apm.WithServerURL("http://192.168.56.114:8200"),
		apm.WithSecretToken("thO5a5ISLcoTrogIcH8XljEPRLs9uqoswl"),
	)
	
	if err != nil {
		log.Printf("APM initialization failed: %v", err)
	} else {
		log.Println("APM initialized successfully for Teacher Service")
	}
}

// APM middleware
func apmMiddleware(handler http.HandlerFunc, operationName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Start transaction
		tx := apm.DefaultTracer().StartTransaction(operationName, "request")
		defer tx.End()
		
		ctx := apm.ContextWithTransaction(r.Context(), tx)
		req := r.WithContext(ctx)
		
		// Add context to request
		handler(w, req)
	}
}

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, traceparent, tracestate")
}

func main() {
	// Initialize APM
	initAPM()

	// Database connection
	if err := database.Connect(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Teacher Routes with APM middleware
	http.HandleFunc("/add-teacher", func(w http.ResponseWriter, r *http.Request) {
		enableCors(w)
		if r.Method == http.MethodOptions {
			return
		}
		apmMiddleware(handlers.AddTeacher, "POST /add-teacher")(w, r)
	})

	http.HandleFunc("/teachers", func(w http.ResponseWriter, r *http.Request) {
		enableCors(w)
		if r.Method == http.MethodOptions {
			return
		}
		apmMiddleware(handlers.GetTeachers, "GET /teachers")(w, r)
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
		apmMiddleware(handlers.DeleteTeacher, "DELETE /delete-teacher")(w, r)
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
		apmMiddleware(handlers.UpdateTeacher, "PUT /update-teacher")(w, r)
	})

	log.Println("Teacher Service running on port 5002")
	log.Fatal(http.ListenAndServe(":5002", nil))
}