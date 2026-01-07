


pipeline {
    agent any
    environment {
        DOCKER_HUB_REPO = "nanil0034/kindergarten-registry"
        SSH_CREDENTIALS = 'vm-ssh-key'
        REMOTE_SERVER_IP = '192.168.121.188'
        GIT_BRANCH = "main"
        SONARQUBE_ENV = "Sonar"
    }

    stages {
        stage('Git pull') {
            steps {
                echo "========Code Cloning========"
                sh """
                    rm -rf Linux_Docker_k8s_Deployment
                    git clone https://github.com/Nabil720/Linux_Docker_k8s_Deployment.git
                    cd Linux_Docker_k8s_Deployment
                """
                echo "Code cloning Successful"
                sh "pwd && ls -la"
            }
        }



        /* ---------------------- Build Multiple Docker Images ---------------------- */

        stage('Build Student Service Docker Image') {
            steps {
                echo "========Building Student Service Docker Image========"
                sh """
                    cd Linux_Docker_k8s_Deployment/kindergarten-registry_k8s/studentservice
                    pwd
                    ls -la
                    docker build -t ${DOCKER_HUB_REPO}-student:${BUILD_NUMBER} .
                """
            }
        }

        stage('Build Teacher Service Docker Image') {
            steps {
                echo "========Building Teacher Service Docker Image========"
                sh """
                    cd Linux_Docker_k8s_Deployment/kindergarten-registry_k8s/teacherservice
                    pwd
                    ls -la
                    docker build -t ${DOCKER_HUB_REPO}-teacher:${BUILD_NUMBER} .
                """
            }
        }

        stage('Build Employee Service Docker Image') {
            steps {
                echo "========Building Employee Service Docker Image========"
                sh """
                    cd Linux_Docker_k8s_Deployment/kindergarten-registry_k8s/employeeservice
                    pwd
                    ls -la
                    docker build -t ${DOCKER_HUB_REPO}-employee:${BUILD_NUMBER} .
                """
            }
        }

        stage('Build Frontend Docker Image') {
            steps {
                echo "========Building Frontend Docker Image========"
                sh """
                    cd Linux_Docker_k8s_Deployment/kindergarten-registry_k8s/frontend
                    pwd
                    ls -la
                    docker build -t ${DOCKER_HUB_REPO}-frontend:${BUILD_NUMBER} .
                """
            }
        }

        stage('Login to Docker Hub & Push Images') {
            steps {
                echo "========Pushing Docker Images to Docker Hub========"
                withCredentials([usernamePassword(credentialsId: 'dockerHubCred', passwordVariable: 'DOCKER_PASS', usernameVariable: 'DOCKER_USER')]) {
                    sh """
                    echo "\$DOCKER_PASS" | docker login -u "\$DOCKER_USER" --password-stdin
                    docker push ${DOCKER_HUB_REPO}-student:${BUILD_NUMBER}
                    docker push ${DOCKER_HUB_REPO}-teacher:${BUILD_NUMBER}
                    docker push ${DOCKER_HUB_REPO}-employee:${BUILD_NUMBER}
                    docker push ${DOCKER_HUB_REPO}-frontend:${BUILD_NUMBER}
                    """
                }
            }
        }

        stage('Update GitOps Repo with New Image Tags') {
            steps {
                echo "========Updating GitOps Repo with New Image Tags========"
                withCredentials([string(credentialsId: 'github-token', variable: 'GITHUB_TOKEN')]) {
                    script {
                        sh '''
                        if [ -d "Linux_Docker_k8s_Deployment" ]; then
                            echo "Directory exists, removing..."
                            rm -rf Linux_Docker_k8s_Deployment
                        fi
                        
                        # Clone and work with main branch instead of main
                        git clone https://x-access-token:${GITHUB_TOKEN}@github.com/Nabil720/Kindergarten-registry_GitOps.git
                        cd Kindergarten-registry_GitOps
                        
                        # Check current branches and switch to main
                        echo "Available branches:"
                        git branch -a
                        
                        # Use main branch (your repository has main, not main)
                        git checkout main

                        # Configure git
                        git config --global user.email "nabilfaruk6@gmail.com"
                        git config --global user.name "Nabil"

                        # Update image tags
                        sed -i "s|image: ${DOCKER_HUB_REPO}-student:.*|image: ${DOCKER_HUB_REPO}-student:${BUILD_NUMBER}|" student-service-deployment.yaml
                        sed -i "s|image: ${DOCKER_HUB_REPO}-teacher:.*|image: ${DOCKER_HUB_REPO}-teacher:${BUILD_NUMBER}|" teacher-service-deployment.yaml
                        sed -i "s|image: ${DOCKER_HUB_REPO}-employee:.*|image: ${DOCKER_HUB_REPO}-employee:${BUILD_NUMBER}|" employee-service-deployment.yaml
                        sed -i "s|image: ${DOCKER_HUB_REPO}-frontend:.*|image: ${DOCKER_HUB_REPO}-frontend:${BUILD_NUMBER}|" frontend-deployment.yaml

                        git add .
                        if git diff --cached --quiet; then
                            echo "No changes to commit."
                        else
                            git commit -m "Jenkins: Update all images to version ${BUILD_NUMBER}"
                            git push origin main
                        fi
                        '''
                    }
                }
            }
        }


        // stage('Deploy to Kubernetes') {
        //     steps {
        //         echo "========Deploying to Kubernetes========"
        //         withCredentials([sshUserPrivateKey(credentialsId: "${SSH_CREDENTIALS}", keyFileVariable: 'SSH_KEY', usernameVariable: 'SSH_USER')]) {
        //             sh """
        //             ssh -o StrictHostKeyChecking=no -i ${SSH_KEY} ${SSH_USER}@${REMOTE_SERVER_IP} "
        //                 cd Linux_Docker_k8s_Deployment/kindergarten-registry_GitOps &&
        //                 kubectl apply -f mongo-deployment.yaml &&
        //                 kubectl apply -f student-service-deployment.yaml &&
        //                 kubectl apply -f teacher-service-deployment.yaml &&
        //                 kubectl apply -f employee-service-deployment.yaml &&
        //                 kubectl apply -f frontend-deployment.yaml &&
        //                 kubectl apply -f ingress.yaml &&
        //                 echo 'Deployment completed successfully!'
        //             "
        //             """
        //         }
        //     }
        // }

    }

    post {
        always {
            echo "Pipeline finished. Collecting artifacts if any."
            sh "docker logout"
        }
        success {
            echo "Pipeline completed successfully."
        }
        failure {
            echo "Pipeline failed. Check logs."
        }
    }
}
