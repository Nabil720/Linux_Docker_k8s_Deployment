package main

import (
    "log"
    "net/http"
    "os"
    "employeeservice/config"
    "employeeservice/database"
    "employeeservice/handlers"

    "go.elastic.co/apm/v2"
)

func initAPM(apmConfig config.APMConfig) {
    // Environment variables set করো APM-এর জন্য
    os.Setenv("ELASTIC_APM_SERVER_URL", apmConfig.ServerURL)
    os.Setenv("ELASTIC_APM_SECRET_TOKEN", apmConfig.SecretToken)
    os.Setenv("ELASTIC_APM_ENVIRONMENT", apmConfig.Environment)
    os.Setenv("ELASTIC_APM_SERVICE_NAME", "employee-service")

    if apm.DefaultTracer().Active() {
        log.Println("APM initialized for Employee Service")
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
    log.Println("Initializing Employee Service with Vault integration...")
    
    vaultClient, err := config.InitVaultClient()
    if err != nil {
        log.Fatalf("Failed to initialize Vault client: %v", err)
    }

    // Vault connection health check
    if err := config.CheckVaultHealth(vaultClient); err != nil {
        log.Printf("Warning: Vault health check: %v", err)
    }

    secrets, err := config.GetSecrets(vaultClient, "employee")
    if err != nil {
        log.Fatalf("Failed to get secrets from Vault: %v", err)
    }

    log.Println("Successfully loaded configuration from Vault")
    log.Printf("Service port: %d", secrets.Port)

    // Step 2: APM initialize করো
    initAPM(secrets.APM)

    // Step 3: Database connection (এখন Vault থেকে URI পাবে)
    os.Setenv("MONGODB_URI", secrets.MongoDB.URI)
    os.Setenv("DATABASE_NAME", secrets.MongoDB.Database)

    log.Printf("Connecting to MongoDB with URI: %s", secrets.MongoDB.URI)
    
    if err := database.Connect(); err != nil {
        log.Fatal("Database connection failed:", err)
    }

    // Step 4: HTTP routes setup করো
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

    // Health check endpoint
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        enableCors(w)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status": "healthy", "service": "employee"}`))
    })

    log.Printf("Employee Service running on port %d", secrets.Port)
    log.Fatal(http.ListenAndServe(":5003", nil))
}