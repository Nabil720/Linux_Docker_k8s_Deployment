package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"employeeservice/database"
	"employeeservice/handlers"

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
	
	if serviceName, ok := secret.Data["elastic-apm-service-name-employee"].(string); ok {
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
		log.Println("APM initialized successfully for Employee Service")
	} else {
		log.Println("APM is not active - check environment variables")
	}
}

// SIMPLIFIED APM middleware - FIXED VERSION
func apmMiddleware(handler http.HandlerFunc, operationName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if APM tracer is available before using it
		tracer := apm.DefaultTracer()
		if tracer == nil || !tracer.Active() {
			// If APM is not available, just call the handler directly
			handler(w, r)
			return
		}
		
		// Start transaction - CORRECT WAY
		tx := tracer.StartTransaction(operationName, "request")
		defer tx.End()
		
		// Set transaction context
		ctx := apm.ContextWithTransaction(r.Context(), tx)
		req := r.WithContext(ctx)
		
		// Call the handler
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
	http.HandleFunc("/emp/health", func(w http.ResponseWriter, r *http.Request) {
		enableCors(w)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "healthy",
			"service": "employee-service",
		})
	})

	// Employee Routes
	http.HandleFunc("/emp/add-employee", func(w http.ResponseWriter, r *http.Request) {
		enableCors(w)
		if r.Method == http.MethodOptions {
			return
		}
		apmMiddleware(handlers.AddEmployee, "POST /add-employee")(w, r)
	})

	http.HandleFunc("/emp/employees", func(w http.ResponseWriter, r *http.Request) {
		enableCors(w)
		if r.Method == http.MethodOptions {
			return
		}
		apmMiddleware(handlers.GetEmployees, "GET /employees")(w, r)
	})

	http.HandleFunc("/emp/delete-employee", func(w http.ResponseWriter, r *http.Request) {
		enableCors(w)
		if r.Method == http.MethodOptions {
			return
		}
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		apmMiddleware(handlers.DeleteEmployee, "DELETE /delete-employee")(w, r)
	})

	http.HandleFunc("/emp/update-employee", func(w http.ResponseWriter, r *http.Request) {
		enableCors(w)
		if r.Method == http.MethodOptions {
			return
		}
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		apmMiddleware(handlers.UpdateEmployee, "PUT /update-employee")(w, r)
	})

	log.Println("Employee Service running on port 5003")
	log.Fatal(http.ListenAndServe(":5003", nil))
}