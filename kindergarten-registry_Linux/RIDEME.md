<h1 style="color: #1bf89cff; font-size: 48px; font-weight: bold;">Linux Based Hosting</h1>

## Kindergarten Registry — React + Go + MongoDB (Deployed with Nginx)

This project is a Kindergarten Registry Management System built using React (frontend), Go (backend), and MongoDB (database).  
It is deployed on an AWS EC2 (Ubuntu t3.micro) instance using Nginx as a reverse proxy.

## Step-by-Step Deployment Guide

### Launch EC2 Instance

1. Create an EC2 instance using **Ubuntu 22.04 (t3.micro)**.
2. Connect via SSH:

```bash
ssh -i key.pem ubuntu@<EC2_PUBLIC_IP>
```

### Clone the Repository

```
cd ~
https://github.com/Nabil720/Linux_Docker_k8s_Deployment.git
```

---

## Project Structure

```
.
├── employeeservice
│   ├── database
│   │   └── db.go
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── handlers
│   │   └── employee.go
│   ├── main.go
│   └── models
│       └── employee.go
├── frontend
│   ├── package.json
│   ├── package-lock.json
│   ├── public
│   │   ├── favicon.ico
│   │   ├── index.html
│   │   ├── logo192.png
│   │   ├── logo512.png
│   │   ├── manifest.json
│   │   └── robots.txt
│   ├── README.md
│   └── src
│       ├── App.css
│       ├── App.js
│       ├── App.test.js
│       ├── components
│       │   ├── EmployeeForm.js
│       │   ├── EmployeeList.js
│       │   ├── StudentForm.js
│       │   ├── StudentList.js
│       │   ├── TeacherForm.js
│       │   └── TeacherList.js
│       ├── index.css
│       ├── index.js
│       ├── logo.svg
│       ├── reportWebVitals.js
│       └── setupTests.js
├── RIDEME.md
├── studentservice
│   ├── database
│   │   └── db.go
│   ├── go.mod
│   ├── go.sum
│   ├── handlers
│   │   └── student.go
│   ├── main.go
│   └── models
│       └── student.go
└── teacherservice
    ├── database
    │   └── db.go
    ├── go.mod
    ├── go.sum
    ├── handlers
    │   └── teacher.go
    ├── main.go
    └── models
        └── teacher.go

```

---

### Install Required Packages

```
sudo apt update && sudo apt upgrade -y
sudo apt install -y nginx docker.io docker-compose nodejs npm golang-go

# Verify installations
nginx -v
docker -v
go version
node -v
```

### Set Up MongoDB

```
cd kindergarten-registry_Linux
nano docker-compose.yml

version: '3.8'

services:
  mongo:
    image: mongo:latest
    container_name: mongo
    restart: unless-stopped
    ports:
      - "27017:27017"
    volumes:
      - mongo_data:/data/db
    environment:
      - MONGO_INITDB_DATABASE=kindergarten

volumes:
  mongo_data:

# Run MongoDB
docker-compose up -d

```

### Setup and Run Backend and Frontend

```
# Student Service (Port 5001):
cd studentservice
go mod tidy
go build
./studentservice &

# Teacher Service (Port 5002):
cd teacherservice
go mod tidy
go build
./teacherservice &

# Employee Service (Port 5003):

cd employeeservice
go mod tidy
go build
./employeeservice &

# Setup and Build Frontend:

cd ../frontend
npm install
npm run build
```

### Configure Nginx

```
# Backup default configuration
sudo mv /etc/nginx/sites-available/default /etc/nginx/sites-available/default.bak

# Create Nginx configuration
sudo nano /etc/nginx/sites-available/kindergarten

server {
    listen 80;
    server_name 192.168.56.112 localhost;

    root /usr/share/nginx/html;
    index index.html index.htm;

    # Frontend routes
    location / {
        try_files $uri $uri/ /index.html;
    }

    # API proxies - Student Service
    location /api/students/ {
        proxy_pass http://localhost:5001/;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }

    # API proxies - Teacher Service
    location /api/teachers/ {
        proxy_pass http://localhost:5002/;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }

    # API proxies - Employee Service
    location /api/employees/ {
        proxy_pass http://localhost:5003/;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}


# Enable configuration
sudo ln -s /etc/nginx/sites-available/kindergarten /etc/nginx/sites-enabled/
sudo rm -f /etc/nginx/sites-enabled/default

# Test and restart Nginx
sudo nginx -t
sudo systemctl restart nginx
sudo systemctl enable nginx
```


### Deploy Frontend Build

```bash
# Copy build files to Nginx directory
sudo cp -r frontend/build/* /usr/share/nginx/html/

# Fix permissions
sudo chown -R www-data:www-data /usr/share/nginx/html/
sudo chmod -R 755 /usr/share/nginx/html/

```



### Access the Application

```
http://192.168.56.112

```

## Applicatio architecture

```
                +---------------------+
                |     Web Browser     |
                | (React Frontend UI) |
                +---------+-----------+
                          |
                          v
                +---------------------+
                |       Nginx         |
                | (Reverse Proxy)     |
                +----+----------+-----+
                    |          |
                    |          |  API Requests
                    |          |
                    v          v
                +-----------+  +-------------------+
                | React     |  |   Go Microservices|
                | Build     |  | +---------------+ |
                | Files     |  | | Student Service| |
                | (static)  |  | |    (5001)     | |
                +-----------+  | +---------------+ |
                              | | Teacher Service| |
                              | |    (5002)     | |
                              | +---------------+ |
                              | | Employee Service|
                              | |    (5003)     | |
                              | +---------------+ |
                              +-------------------+
                                        |
                                        v
                              +-------------------+
                              |     MongoDB       |
                              | (Docker Container)|
                              |   localhost:27017 |
                              +-------------------+
                  



```

![Website View](./Images/Screenshot%20from%202025-11-05%2018-09-16.png)
![Website View](./Images/Screenshot%20from%202025-11-05%2018-09-28.png)
![Website View](./Images/Screenshot%20from%202025-11-05%2018-14-00.png)
![Website View](./Images/Screenshot%20from%202025-11-05%2018-09-55.png)
![Website View](./Images/Screenshot%20from%202025-11-05%2018-10-03.png)
![Website View](./Images/Screenshot%20from%202025-11-05%2018-11-54.png)




