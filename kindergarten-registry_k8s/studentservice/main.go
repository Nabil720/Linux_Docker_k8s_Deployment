package main

import (
	"log"
	"net/http"
	"os"
	"studentservice/config"
	"studentservice/database"
	"studentservice/handlers"

	"go.elastic.co/apm/v2"
)

func initAPM(apmConfig config.APMConfig) {
	// Environment variables set করো APM-এর জন্য
	os.Setenv("ELASTIC_APM_SERVER_URL", apmConfig.ServerURL)
	os.Setenv("ELASTIC_APM_SECRET_TOKEN", apmConfig.SecretToken)
	os.Setenv("ELASTIC_APM_ENVIRONMENT", apmConfig.Environment)
	os.Setenv("ELASTIC_APM_SERVICE_NAME", "student-service")

	if apm.DefaultTracer().Active() {
		log.Println("APM initialized for Student Service")
	} else {
		log.Println("APM not active - check Vault configuration")
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
	// Step 1: Vault থেকে secrets লোড করো
	vaultClient, err := config.InitVaultClient()
	if err != nil {
		log.Fatalf("Failed to initialize Vault client: %v", err)
	}

	secrets, err := config.GetSecrets(vaultClient, "student")
	if err != nil {
		log.Fatalf("Failed to get secrets from Vault: %v", err)
	}

	log.Println("Successfully loaded configuration from Vault")

	// Step 2: APM initialize করো
	initAPM(secrets.APM)

	// Step 3: Database connection (এখন Vault থেকে URI পাবে)
	os.Setenv("MONGODB_URI", secrets.MongoDB.URI)
	os.Setenv("DATABASE_NAME", secrets.MongoDB.Database)

	if err := database.Connect(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Step 4: HTTP routes setup করো
	http.HandleFunc("/std/add-student", func(w http.ResponseWriter, r *http.Request) {
		enableCors(w)
		if r.Method == http.MethodOptions {
			return
		}
		apmMiddleware(handlers.AddStudent, "POST /add-student")(w, r)
	})

	http.HandleFunc("/std/students", func(w http.ResponseWriter, r *http.Request) {
		enableCors(w)
		if r.Method == http.MethodOptions {
			return
		}
		apmMiddleware(handlers.GetStudents, "GET /students")(w, r)
	})

	http.HandleFunc("/std/delete-student", func(w http.ResponseWriter, r *http.Request) {
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

	http.HandleFunc("/std/update-student", func(w http.ResponseWriter, r *http.Request) {
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

	log.Printf("Student Service running on port %d", secrets.Port)
	log.Fatal(http.ListenAndServe(":5001", nil))
}