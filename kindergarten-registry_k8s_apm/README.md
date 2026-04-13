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
---

## Step 1: Prepare Vault Server

### Terminal 1 (Vault Server):

1. **Set the Vault address**  
   Set the Vault server address (you may need to adjust the address based on your configuration):
   ```bash
   export VAULT_ADDR='http://localhost:8200'
   ```
 2. **Login to Vault**
    Use your Vault token to log in:
    ```bash
    vault login **
    ```
 3. **Check Vault status**
    Verify that the Vault server is running and healthy:
    ```bash
    vault status
    ```
4. **Enable the secret engine (if not enabled)**
   If the kv-v2 secret engine is not enabled, enable it:
   ```bash
   vault secrets list | grep secret/ || vault secrets enable -path=kindergarten/config kv
   ```
5. **Create a test secret**
   Create a test secret that will later be injected into the Kubernetes pod:
   ```bash
   # MongoDB credentials
   vault kv put kindergarten/config/mongodb \
   username="myUser" \
   password="myPassword" \
   database="kindergarten" \
   uri="mongodb://myUser:myPassword@mongo:27017/kindergarten?authSource=admin"

   # Elastic APM configuration
   vault kv put kindergarten/config/apm \
   server_url="http://192.168.121.224:8200" \
   secret_token="your_apm_access_token" \
   environment="production"

   # Service-specific APM names
   vault kv put kindergarten/config/services \
   student="student-service" \
   teacher="teacher-service" \
   employee="employee-service"

   # Service port configurations
   vault kv put kindergarten/config/ports \
   student=5001 \
   teacher=5002 \
   employee=5003

   ```
   6. **Verify the secret**
      Retrieve and verify the secret:
   ```bash
   vault kv get kindergarten/config/mongodb
   vault kv get kindergarten/config/apm
   vault kv get kindergarten/config/services
   vault kv get kindergarten/config/services
   vault kv get kindergarten/config/ports
   ```

## Step 2: Install Vault Agent Injector

### Terminal 2 (Kubernetes Master):

1. **Add the Helm repository**
   Add the HashiCorp Helm repository:
   ```bash
   helm repo add hashicorp https://helm.releases.hashicorp.com
   helm repo update
   ```
2. **Install Vault Agent Injector**
   Install the Vault Agent Injector into Kubernetes (Note: Vault server is external in this setup):
   ```bash
   helm install vault hashicorp/vault \
   --set "injector.enabled=true" \
   --set "server.enabled=false" \
   --set "global.externalVaultAddr=http://192.168.61.164:8200" \  # vault server IP or HA-proxy ip (If we setup  vault as cluster)
   --namespace vault \
   --create-namespace
   ```
3. **Check installation**
   ```bash
   kubectl get pods -n vault
   kubectl get svc -n vault
   ```
## Step 3: Collect Kubernetes Configuration

### Terminal 2 (Kubernetes Master):

1. **Get Kubernetes API URL**
   ```bash
   K8S_API=$(kubectl config view --minify -o jsonpath='{.clusters[0].cluster.server}')
   echo "$K8S_API"
   ```
2. **Get proper CA Certificate**
   ```bash
   echo "2. CA_CERTIFICATE (proper format):"
   if [ -f "/etc/kubernetes/pki/ca.crt" ]; then
   sudo cat /etc/kubernetes/pki/ca.crt
   else
   kubectl config view --raw -o jsonpath='{.clusters[0].cluster.certificate-authority-data}' | base64 --decode
   fi
   echo ""
   ```
3. **Create Vault service account**
   ```bash
   kubectl create serviceaccount vault-auth 2>/dev/null || echo "Service account already exists"
   ```
4. **Create cluster role binding**
   ```bash
   kubectl create clusterrolebinding vault-auth-binding \
   --clusterrole=system:auth-delegator \
   --serviceaccount=default:vault-auth 2>/dev/null || echo "Cluster role binding already exists"
   ```
5. **Get  issuer**
   ```bash
   kubectl get --raw /.well-known/openid-configuration
   ```
## Step 4: Configure Vault Kubernetes Auth Method

### Terminal 1 (Vault Server):

1. **Enable Kubernetes authentication method**
   Enable the Kubernetes authentication method in Vault:
   ```bash
   vault auth enable kubernetes
   ```

2. **Configure Vault with Kubernetes details**
   ```bash
   KUBERNETES_HOST="https://10.70.57.50:6443"  # Your Kubernetes API URL
   CA_CERT="<CA_CERTIFICATE>"  # Kubernetes CA certificate

   vault write auth/kubernetes/config \
     kubernetes_host="$KUBERNETES_HOST" \
     kubernetes_ca_cert="$CA_CERT" \
     issuer="https://kubernetes.default.svc.cluster.local" \
     disable_iss_validation=false
   ```
3. **Create a Vault policy**
   Create a policy that allows read access to the secret:
   ```bash
   vault policy write kindergarten-policy - <<EOF
   path "kindergarten/config/*" {
   capabilities = ["create", "update", "read"]
   }
   EOF
   ```
4. **Create a role for Kubernetes authentication**
   Bind the myapp-policy to the role and link it with the service account and namespace:
   ```bash
   vault write auth/kubernetes/role/kindergarten-role \
   bound_service_account_names=vault-auth \
   bound_service_account_namespaces=default \
   policies=kindergarten-policy \
   ttl=24
   ```
5. **Verify the configuration**
   Check the configuration:
   ```bash
   vault read auth/kubernetes/config
   vault read auth/kubernetes/role/kindergarten-role
   ```
## Step 5: Create and Apply Test Pod

## Terminal 2 (Kubernetes Master):

1. **Create a test pod YAML file**
   Create a simple pod YAML file to test the secret injection:
   ```bash
   cat > teacher-service-deployment.yaml <<EOF
    apiVersion: apps/v1
    kind: Deployment
    metadata:
    name: teacher-service
    labels:
        app: teacher-service
    spec:
    replicas: 1
    selector:
        matchLabels:
        app: teacher-service
    template:
        metadata:
        labels:
            app: teacher-service
        annotations:
            vault.hashicorp.com/agent-inject: "true"
            vault.hashicorp.com/tls-skip-verify: "true"
            vault.hashicorp.com/agent-pre-populate-only: "true"
            vault.hashicorp.com/agent-inject-status: "update"
            vault.hashicorp.com/role: "kindergarten-role"

            # Inject all secrets into a single env file
            vault.hashicorp.com/agent-inject-secret-env: "kindergarten/config/all-secrets"
            vault.hashicorp.com/agent-inject-template-env: |
            {{- with secret "kindergarten/config/mongodb" -}}
            {{- range $k, $v := .Data }}
            {{ $k }}={{ $v }}
            {{- end }}
            {{- end }}
            {{- with secret "kindergarten/config/ports" -}}
            {{- range $k, $v := .Data }}
            {{ $k }}={{ $v }}
            {{- end }}
            {{- end }}
            {{- with secret "kindergarten/config/services" -}}
            {{- range $k, $v := .Data }}
            {{ $k }}={{ $v }}
            {{- end }}
            {{- end }}
        spec:
        serviceAccountName: vault-auth
        automountServiceAccountToken: true
        containers:
            - name: teacher-service
            image: nanil0034/kindergarten-registry-teacher:112
            imagePullPolicy: Always
            ports:
                - containerPort: 5002
            command: ["/bin/sh"]
            args:
                - "-c"
                - |
                # Export all secrets from the single env file
                if [ -f /vault/secrets/env ]; then
                    export $(cat /vault/secrets/env | xargs)
                fi

                # Start the application
                exec /app/main

   ---
   apiVersion: v1
   kind: Service
   metadata:
   name: teacher-service
   spec:
   selector:
      app: teacher-service
   ports:
      - protocol: TCP
         port: 5002
         targetPort: 5002
         nodePort: 30002
   type: NodePort
   EOF

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
