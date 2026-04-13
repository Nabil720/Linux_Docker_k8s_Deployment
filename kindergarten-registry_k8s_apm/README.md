#  Kindergarten Registry — Kubernetes + Vault + APM Deployment

A full-stack **React + Go Microservices + MongoDB** application containerized with Docker, published to Docker Hub, and deployed on a **Kubernetes** cluster with **HashiCorp Vault** secrets management and **Elastic APM** observability integration.

---

##  Table of Contents

1. [Project Overview](#project-overview)
2. [Architecture](#architecture)
3. [Project Structure](#project-structure)
4. [Docker — Build & Push to Docker Hub](#docker--build--push-to-docker-hub)
5. [Kubernetes Deployment](#kubernetes-deployment)
6. [HashiCorp Vault Integration](#hashicorp-vault-integration)
7. [Elastic APM Integration](#elastic-apm-integration)
8. [MongoDB Setup Inside Cluster](#mongodb-setup-inside-cluster)
9. [Screenshots](#screenshots)

---

## Project Overview

The **Kindergarten Registry** system manages **Students**, **Teachers**, and **Employees** through independent Go microservices behind a React frontend. Each service is independently Docker-built, pushed to Docker Hub, and deployed via Kubernetes manifests.

| Technology | Role |
|---|---|
| React (Nginx) | Frontend UI |
| Go | Microservices (Student / Teacher / Employee) |
| MongoDB | Persistent database |
| Docker Hub | Image registry |
| Kubernetes | Container orchestration |
| HashiCorp Vault | Secrets management |
| Elastic Stack (ELK + APM) | Observability & monitoring |

---

## Architecture

```
                    +--------------------+
                    |    Web Browser     |
                    |  (React Frontend)  |
                    +--------+-----------+
                             |
                        LoadBalancer
                             |
                    +--------+-----------+
                    |   Nginx (port 80)  |
                    |  frontend Service  |
                    +--+------+-------+--+
                       |      |       |
              /student  |  /teacher   | /employee
                       v      v       v
           +----------+ +----------+ +----------+
           |  Student | | Teacher  | | Employee |
           | Service  | | Service  | | Service  |
           | :5001    | | :5002    | | :5003    |
           +----+-----+ +----+-----+ +----+-----+
                |             |            |
                +------+-------+------+----+
                              |
                    +---------+----------+
                    |  MongoDB StatefulSet|
                    |  (mongo:6.0 :27017) |
                    |  PVC: 1Gi          |
                    +--------------------+







                   +----------------------+
                   |   HashiCorp Vault    |
                   |----------------------|
                   | VAULT_SKIP_VERIFY=true
                   | MONGODB_URI          |
                   | DATABASE_NAME        |
                   +----------+-----------+
                              |
                              |  (Env Variables Injected)
                              v
        +----------------+----------------+----------------+
        |   Student      |    Teacher     |    Employee     |
        |   Services     |    Services    |    Services     |
        +--------+-------+--------+-------+--------+--------+
                 \                |                /
                  \               |               /
                   \              |              /
                    \             |             /
                     +------------+------------+
                                  |
                                  |   APM Traces
                                  v
                     +----------------------------+
                     |        APM Stack           |
                     |----------------------------|
                     | APM Server      :8200      |
                     | Elasticsearch  :9200       |
                     | Kibana         :5601       |
                     +----------------------------+
```

---

## Project Structure

```
kindergarten-registry_k8s_apm/
├── frontend/                          # React app
│   ├── Dockerfile
│   └── src/
│       └── components/
│           ├── StudentForm.js / StudentList.js
│           ├── TeacherForm.js / TeacherList.js
│           └── EmployeeForm.js / EmployeeList.js
├── studentservice/                    # Go microservice — port 5001
│   ├── Dockerfile
│   ├── main.go
│   ├── handlers/student.go
│   ├── models/student.go
│   └── database/db.go
├── teacherservice/                    # Go microservice — port 5002
│   ├── Dockerfile
│   ├── main.go
│   ├── handlers/teacher.go
│   ├── models/teacher.go
│   └── database/db.go
├── employeeservice/                   # Go microservice — port 5003
│   ├── Dockerfile
│   ├── main.go
│   ├── handlers/employee.go
│   ├── models/employee.go
│   └── database/db.go
├── k8s_manifest/
│   ├── frontend-deployment.yaml
│   ├── student-service-deployment.yaml
│   ├── teacher-service-deployment.yaml
│   ├── employee-service-deployment.yaml
│   └── mongo-deployment.yaml
├── mongo-init.js
└── Images/                            # Screenshots
```

---

## Docker — Build & Push to Docker Hub

Each service has its own `Dockerfile`. Images are built and pushed to **Docker Hub** under the `nanil0034` namespace.

### Build & Push — Frontend

```bash
cd frontend/
docker build -t nanil0034/kindergarten-registry-frontend:latest .
docker push nanil0034/kindergarten-registry-frontend:latest
```

### Build & Push — Student Service

```bash
cd studentservice/
docker build -t nanil0034/kindergarten-registry-student:latest .
docker push nanil0034/kindergarten-registry-student:latest
```

### Build & Push — Teacher Service

```bash
cd teacherservice/
docker build -t nanil0034/kindergarten-registry-teacher:latest .
docker push nanil0034/kindergarten-registry-teacher:latest
```

### Build & Push — Employee Service

```bash
cd employeeservice/
docker build -t nanil0034/kindergarten-registry-employee:latest .
docker push nanil0034/kindergarten-registry-employee:latest
```

### Docker Hub Images

| Image | Tag |
|---|---|
| `nanil0034/kindergarten-registry-frontend` | `latest` |
| `nanil0034/kindergarten-registry-student` | `latest` |
| `nanil0034/kindergarten-registry-teacher` | `latest` |
| `nanil0034/kindergarten-registry-employee` | `latest` |

---


## APM Installation Guide

```
https://github.com/siyamsarker/elastic-apm-quickstart
```
## Vault Installation Guide

```
https://github.com/Nabil720/Hashicorp-Vault/tree/master/Vault_Apt  # Here the installation  is three node cluster , You can use standalone also
```

## Kubernetes Deployment

All manifests live in `k8s_manifest/`. Apply them to your cluster:

```bash
# Apply all manifests at once
kubectl apply -f k8s_manifest/

# Or apply individually
kubectl apply -f k8s_manifest/mongo-deployment.yaml
kubectl apply -f k8s_manifest/student-service-deployment.yaml
kubectl apply -f k8s_manifest/teacher-service-deployment.yaml
kubectl apply -f k8s_manifest/employee-service-deployment.yaml
kubectl apply -f k8s_manifest/frontend-deployment.yaml
```

### Verify Deployments

```bash
kubectl get pods
kubectl get services
kubectl get pv,pvc
```

### Services & Ports

| Service | Type | Port |
|---|---|---|
| `frontend` | LoadBalancer | 80 |
| `student-service` | ClusterIP | 5001 |
| `teacher-service` | ClusterIP | 5002 |
| `employee-service` | ClusterIP | 5003 |
| `mongo` | ClusterIP | 27017 |


### Access the Frontend

```bash
# Get the LoadBalancer external IP
kubectl get service frontend

# Or port-forward for local testing
kubectl port-forward service/frontend 8080:80
# → http://localhost:8080
```

---

## HashiCorp Vault Integration

**HashiCorp Vault** injects secrets into the Go microservices as **environment variables** at pod startup. The Go code reads them using `os.Getenv()` — no Vault API calls are made inside the application itself.

### How It Works

1. Vault stores the MongoDB credentials (`MONGODB_URI`, `DATABASE_NAME`) as secrets.
2. At pod creation, the **Vault Agent Injector** injects those secrets as environment variables into each container.
3. The Kubernetes manifest sets `VAULT_SKIP_VERIFY=true` to allow the Vault injector to work with self-signed TLS certificates.
4. Each Go service reads the injected values at startup:

```
 # Vault Integration
https://github.com/Nabil720/Hashicorp-Vault/blob/master/Vault_injector/README.md
```
---

## Elastic APM Integration

Each Go microservice is fully instrumented with the **Elastic APM Go agent** (`go.elastic.co/apm/v2`). Every HTTP endpoint is wrapped in an `apmMiddleware()` that creates an APM transaction per request — so **all service traffic is visible in Kibana APM** in real time.

### How APM is Integrated in the Code

Each service's `main.go` initializes APM and wraps every route:

```go
// main.go — APM middleware wraps every route
func apmMiddleware(handler http.HandlerFunc, operationName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        tracer := apm.DefaultTracer()
        tx := tracer.StartTransaction(operationName, "request")
        defer tx.End()
        ctx := apm.ContextWithTransaction(r.Context(), tx)
        handler(w, r.WithContext(ctx))
    }
}

// Every route is traced:
http.HandleFunc("/add-student", func(w http.ResponseWriter, r *http.Request) {
    apmMiddleware(handlers.AddStudent, "POST /add-student")(w, r)
})
http.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
    apmMiddleware(handlers.GetStudents, "GET /students")(w, r)
})
// ... same for /delete-student, /update-student
```

### APM Dependency

```go
// go.mod
require (
    go.elastic.co/apm/v2 v2.4.7
)
```

### APM Environment Variables (set per pod)

```bash
ELASTIC_APM_SERVER_URL=http://<APM_SERVER_HOST>:8200
ELASTIC_APM_SECRET_TOKEN=B7n5dCdEDTDppEbm
ELASTIC_APM_SERVICE_NAME=student-service   # teacher-service / employee-service
ELASTIC_APM_ENVIRONMENT=production
```

### Stack Components

| Component | URL | Purpose |
|---|---|---|
| Elasticsearch | `http://localhost:9200` | Stores all APM trace data |
| Kibana | `http://localhost:5601` | Visualize traces & dashboards |
| APM Server | `http://localhost:8200` | Receives traces from Go services |

### Credentials

```
Username:         elastic
Password:         HEtL6W7qxEUJcs20
APM Secret Token: B7n5dCdEDTDppEbm
```

### Viewing All Service Traffic in Kibana

1. Open Kibana → `http://localhost:5601`
2. Navigate to **Observability → APM**
3. Select a service (`student-service`, `teacher-service`, or `employee-service`)
4. View **Transactions** — every API call (GET /students, POST /add-student, etc.)
5. View **Traces** — full request timelines
6. View **Errors** and **Metrics** — latency, throughput, error rates

### APM Status Check

```bash
curl http://localhost:9200/_cluster/health   # Elasticsearch
curl http://localhost:5601/api/status        # Kibana
curl http://localhost:8200/                  # APM Server
```

```
🔍 Final status check:
  - Elasticsearch: Running
  - Kibana:        Running
  - APM Server:    Running (auth required)
```

---

## MongoDB Setup Inside Cluster

After deploying MongoDB, initialize the database and user:

```bash
# Exec into the MongoDB pod
kubectl exec -it <mongo-pod-name> -- mongosh

# Inside mongosh — switch to admin
use admin;

# Create application user
db.createUser({
  user: "myUser",
  pwd: "myPassword",
  roles: [
    { role: "readWrite", db: "kindergarten" },
    { role: "readWrite", db: "admin" },
    { role: "dbAdmin", db: "kindergarten" },
    { role: "userAdminAnyDatabase", db: "admin" }
  ]
});

# Switch to app database
use kindergarten;

# Create collections
db.createCollection("employees");
db.createCollection("students");
db.createCollection("teachers");

show collections;
exit
```

---

## Screenshots

### Application UI

![Website View](./Images/image_original)

### Kubernetes Cluster — Pods & Services

![K8s Pods](./Images/image_original%20(1))

![K8s Services](./Images/image_original%20(2))

### Elastic APM — Kibana Dashboard

![APM Dashboard](./Images/image_original%20(3))

![APM Traces](./Images/image_original%20(4))

>  Image files sourced from the `Images/` directory in the repository.

---

## 👤 Author

**Nabil** — [GitHub: Nabil720](https://github.com/Nabil720)

> Docker Hub: [nanil0034](https://hub.docker.com/u/nanil0034)
