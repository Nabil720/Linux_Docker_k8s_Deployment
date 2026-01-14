package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"studentservice/database"
	"studentservice/handlers"

	"go.elastic.co/apm/v2"
	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/kubernetes"
)

// ভাউল্ট থেকে secrets লোড করার ফাংশন
func loadEnvFromVault() {
	log.Println("Loading secrets from Vault...")

	// ভাউল্ট কনফিগারেশন
	config := vault.DefaultConfig()
	config.Address = "http://192.168.121.132:8200"

	// ভাউল্ট ক্লায়েন্ট তৈরি
	client, err := vault.NewClient(config)
	if err != nil {
		log.Printf("Warning: Failed to create Vault client: %v", err)
		log.Println("Using default environment variables...")
		return
	}

	// Kubernetes authentication
	k8sAuth, err := auth.NewKubernetesAuth("kindergarten-role")
	if err != nil {
		log.Printf("Warning: Failed to create Kubernetes auth: %v", err)
		return
	}

	// ভাউল্টে লগইন
	authInfo, err := client.Auth().Login(context.Background(), k8sAuth)
	if err != nil {
		log.Printf("Warning: Failed to login to Vault: %v", err)
		log.Println("Using default environment variables...")
		return
	}
	
	if authInfo == nil {
		log.Println("Warning: No auth info received from Vault")
		return
	}

	// Secrets পড়ুন
	secret, err := client.KVv2("kindergarten").Get(context.Background(), "config")
	if err != nil {
		log.Printf("Warning: Failed to read secrets from Vault: %v", err)
		return
	}

	// Secrets থেকে environment variables সেট করুন
	if mongoURI, ok := secret.Data["mongodb-uri"].(string); ok {
		os.Setenv("MONGODB_URI", mongoURI)
		log.Println("MONGODB_URI loaded from Vault")
	}
	
	if dbName, ok := secret.Data["database-name"].(string); ok {
		os.Setenv("DATABASE_NAME", dbName)
		log.Println("DATABASE_NAME loaded from Vault")
	}
	
	if apmURL, ok := secret.Data["elastic-apm-server-url"].(string); ok {
		os.Setenv("ELASTIC_APM_SERVER_URL", apmURL)
		log.Println("ELASTIC_APM_SERVER_URL loaded from Vault")
	}
	
	if apmToken, ok := secret.Data["elastic-apm-secret-token"].(string); ok {
		os.Setenv("ELASTIC_APM_SECRET_TOKEN", apmToken)
		log.Println("ELASTIC_APM_SECRET_TOKEN loaded from Vault")
	}
	
	if serviceName, ok := secret.Data["elastic-apm-service-name-student"].(string); ok {
		os.Setenv("ELASTIC_APM_SERVICE_NAME", serviceName)
		log.Println("ELASTIC_APM_SERVICE_NAME loaded from Vault")
	}

	log.Println("All secrets loaded successfully from Vault!")
}

func initAPM() {
	// প্রথমে ভাউল্ট থেকে secrets লোড করুন
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

	// Health check endpoint
	http.HandleFunc("/std/health", func(w http.ResponseWriter, r *http.Request) {
		enableCors(w)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "healthy",
			"service": "student-service",
		})
	})

	// Student Routes
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

	log.Println("Student Service running on port 5001")
	log.Fatal(http.ListenAndServe(":5001", nil))
}