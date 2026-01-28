package main

import (
    "log"
    "net/http"
    "os"
    "strconv"
    "employeeservice/database"
    "employeeservice/handlers"
)

func enableCors(w http.ResponseWriter) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func main() {
    log.Println("Initializing Employee Service...")
    
    // Step 1: Environment variables 
    mongoURI := os.Getenv("MONGODB_URI")
    if mongoURI == "" {
        log.Println("MONGODB_URI not set in environment")
    }
    
    portStr := os.Getenv("SERVICE_PORT")
    port := 5003
    if portStr != "" {
        if p, err := strconv.Atoi(portStr); err == nil {
            port = p
        } else {
            log.Println("Invalid SERVICE_PORT value, using default port 5003")
        }
    } else {
        log.Println("SERVICE_PORT not set in environment, using default port 5003")
    }
    
    serviceName := os.Getenv("SERVICE_NAME")
    if serviceName == "" {
        log.Println("SERVICE_NAME not set in environment")
    }
    
    // Log the configuration details
    log.Printf("Service: %s", serviceName)
    log.Printf("Port: %d", port)

    // Step 2: Database connection
    if mongoURI != "" {
        os.Setenv("MONGODB_URI", mongoURI)
        if err := database.Connect(); err != nil {
            log.Fatal("Database connection failed:", err)
        }
    } else {
        log.Fatal("Cannot proceed without MongoDB URI")
    }

    // Step 3: HTTP routes setup
    setupRoutes()
    
    log.Printf("Employee Service running on port %d", port)
    log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

func setupRoutes() {
    http.HandleFunc("/emp/add-employee", func(w http.ResponseWriter, r *http.Request) {
        enableCors(w)
        if r.Method == http.MethodOptions {
            return
        }
        handlers.AddEmployee(w, r)
    })

    http.HandleFunc("/emp/employees", func(w http.ResponseWriter, r *http.Request) {
        enableCors(w)
        if r.Method == http.MethodOptions {
            return
        }
        handlers.GetEmployees(w, r)
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
        handlers.DeleteEmployee(w, r)
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
        handlers.UpdateEmployee(w, r)
    })

    // Health check endpoint
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        enableCors(w)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status": "healthy", "service": "employee"}`))
    })
}
