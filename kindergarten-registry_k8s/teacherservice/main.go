package main

import (
	"log"
	"net/http"
	"teacherservice/database" 
	"teacherservice/handlers"

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
	initAPM()
	// Database connection
	if err := database.Connect(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Teacher Routes
	http.HandleFunc("/tech/add-teacher", func(w http.ResponseWriter, r *http.Request) {   // Here using this  "/tech/add-teacher" path , we POST in  the Add-teacher service 
		enableCors(w)
		if r.Method == http.MethodOptions {
			return
		}
		apmMiddleware(handlers.AddTeacher, "POST /add-teacher")(w, r)
	})

	http.HandleFunc("/tech/teachers", func(w http.ResponseWriter, r *http.Request) {      // Here using this  "/tech/teachers" path , we GET teachers  
		enableCors(w)
		enableCors(w)
		if r.Method == http.MethodOptions {
			return
		}
		apmMiddleware(handlers.GetTeachers, "GET /teachers")(w, r)
	})

	http.HandleFunc("/tech/delete-teacher", func(w http.ResponseWriter, r *http.Request) {    // Here using this  "/tech/delete-teacher" path , we DELETE in  the Teachers   service 
		enableCors(w)
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

	http.HandleFunc("/tech/update-teacher", func(w http.ResponseWriter, r *http.Request) {   // Here using this  "/tech/update-teacher" path , we PUT in  the Teachers service 
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