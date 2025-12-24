<h1 style="color: #1bf89cff; font-size: 48px; font-weight: bold;">Docker based hosting</h1>


## Kindergarten Registry

A full-stack microservices application built with React frontend, polyglot backend services (Go, Python, Spring Boot), and MongoDB, deployed using Docker Compose.


## Step-by-Step Deployment Guide

### Launch EC2 Instance
1. Create an EC2 instance using **Ubuntu 22.04 (t3.micro)**.  
2.  Connect via SSH:
   ```bash
   ssh -i key.pem ubuntu@<EC2_PUBLIC_IP>
```

3. Clone the repository

```bash
cd ~
https://github.com/Nabil720/Linux_Docker_k8s_Deployment.git
```
4. Update system and install dependencies

```
sudo apt update -y
sudo apt upgrade -y
sudo apt install docker.io docker-compose -y
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker $USER
newgrp docker
```



```
Linux_Docker_k8s_Deployment/kindergarten-registry_Docker/
├── docker-compose.yml
├── employeeservice-python/         
│   ├── app/
│   │   ├── __init__.py
│   │   ├── models.py
│   │   ├── routes.py
│   │   └── database.py
│   ├── requirements.txt
│   ├── Dockerfile
│   └── run.py
├── frontend
│   ├── Dockerfile
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
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── handlers
│   │   └── student.go
│   ├── main.go
│   └── models
│       └── student.go
├── teacherservice-springboot/      
│   ├── src/
│   │   └── main/
│   │       ├── java/
│   │       │   └── com/
│   │       │       └── kindergarten/
│   │       │           └── teacherservice/
│   │       │               ├── TeacherServiceApplication.java
│   │       │               ├── controller/
│   │       │               │   └── TeacherController.java
│   │       │               ├── model/
│   │       │               │   └── Teacher.java
│   │       │               ├── repository/
│   │       │               │   └── TeacherRepository.java
│   │       │               └── service/
│   │       │                   └── TeacherService.java
│   │       └── resources/
│   │           └── application.properties
│   ├── Dockerfile
│   ├── pom.xml
│   └── mvnw
```


## docker-compose up -d




##  Applicatio architecture


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
                | React     |  |   Microservices   |
                | Build     |  | +---------------+ |
                | Files     |  | | Student Service| |
                | (static)  |  | |    Go (5001)  | |
                +-----------+  | +---------------+ |
                              | | Teacher Service| |
                              | |Spring Boot(5002)| |
                              | +---------------+ |
                              | | Employee Service|
                              | |  Flask (5003)  | |
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


![Website View](./Images/Screenshot%20from%202025-11-05%2018-35-13.png)


![Instance View](./Images/Screenshot%20from%202025-11-05%2018-35-26.png)


![Instance View](./Images/Screenshot%20from%202025-11-05%2018-35-30.png)

![Instance View](./Images/Screenshot%20from%202025-11-05%2018-35-38.png)

