pipeline {
    agent any

    environment {
        DOCKER_HUB_REPO = "nanil0034/kindergarten-registry"
        SSH_CREDENTIALS = 'vm-ssh-key'
        REMOTE_SERVER_IP = '192.168.121.188'
        GIT_BRANCH = "master"
    }

    stages {

        /* ---------------------- Verify SCM Checkout ---------------------- */
        stage('Verify Workspace') {
            steps {
                echo "======== Verifying Jenkins SCM Checkout ========"
                sh '''
                    echo "Current workspace:"
                    pwd
                    echo "Files:"
                    ls -la
                '''
            }
        }

        /* ---------------------- Build Docker Images ---------------------- */

        stage('Build Student Service Image') {
            steps {
                sh '''
                    cd kindergarten-registry_k8s/studentservice
                    docker build -t ${DOCKER_HUB_REPO}-student:${BUILD_NUMBER} .
                '''
            }
        }

        stage('Build Teacher Service Image') {
            steps {
                sh '''
                    cd kindergarten-registry_k8s/teacherservice
                    docker build -t ${DOCKER_HUB_REPO}-teacher:${BUILD_NUMBER} .
                '''
            }
        }

        stage('Build Employee Service Image') {
            steps {
                sh '''
                    cd kindergarten-registry_k8s/employeeservice
                    docker build -t ${DOCKER_HUB_REPO}-employee:${BUILD_NUMBER} .
                '''
            }
        }

        stage('Build Frontend Image') {
            steps {
                sh '''
                    cd kindergarten-registry_k8s/frontend
                    docker build -t ${DOCKER_HUB_REPO}-frontend:${BUILD_NUMBER} .
                '''
            }
        }

        /* ---------------------- Push Images to Docker Hub ---------------------- */

        stage('Docker Login & Push') {
            steps {
                withCredentials([
                    usernamePassword(
                        credentialsId: 'dockerHubCred',
                        usernameVariable: 'DOCKER_USER',
                        passwordVariable: 'DOCKER_PASS'
                    )
                ]) {
                    sh '''
                        echo "$DOCKER_PASS" | docker login -u "$DOCKER_USER" --password-stdin
                        docker push ${DOCKER_HUB_REPO}-student:${BUILD_NUMBER}
                        docker push ${DOCKER_HUB_REPO}-teacher:${BUILD_NUMBER}
                        docker push ${DOCKER_HUB_REPO}-employee:${BUILD_NUMBER}
                        docker push ${DOCKER_HUB_REPO}-frontend:${BUILD_NUMBER}
                    '''
                }
            }
        }

        /* ---------------------- Update GitOps Repository ---------------------- */

        stage('Update GitOps Repo Image Tags') {
            steps {
                withCredentials([
                    string(credentialsId: 'github-token', variable: 'GITHUB_TOKEN')
                ]) {
                    sh '''
                        rm -rf gitops-repo

                        git clone https://x-access-token:${GITHUB_TOKEN}@github.com/Nabil720/Linux_Docker_k8s_Deployment.git gitops-repo
                        cd gitops-repo
                        git checkout master

                        git config user.email "nabilfaruk6@gmail.com"
                        git config user.name "Nabil"

                        cd kindergarten-registry_GitOps

                        sed -i "s|image: ${DOCKER_HUB_REPO}-student:.*|image: ${DOCKER_HUB_REPO}-student:${BUILD_NUMBER}|" student-service-deployment.yaml
                        sed -i "s|image: ${DOCKER_HUB_REPO}-teacher:.*|image: ${DOCKER_HUB_REPO}-teacher:${BUILD_NUMBER}|" teacher-service-deployment.yaml
                        sed -i "s|image: ${DOCKER_HUB_REPO}-employee:.*|image: ${DOCKER_HUB_REPO}-employee:${BUILD_NUMBER}|" employee-service-deployment.yaml
                        sed -i "s|image: ${DOCKER_HUB_REPO}-frontend:.*|image: ${DOCKER_HUB_REPO}-frontend:${BUILD_NUMBER}|" frontend-deployment.yaml

                        git add .
                        git commit -m "Jenkins: update images to build ${BUILD_NUMBER}" || echo "No changes to commit"
                        git push origin master
                    '''
                }
            }
        }

    }

    post {
        always {
            echo "Pipeline finished"
            sh "docker logout || true"
        }
        success {
            echo "Pipeline SUCCESS 🎉"
        }
        failure {
            echo "Pipeline FAILED ❌"
        }
    }
}
