pipeline {
    agent any

    environment {
        IMAGE_NAME = "ttl.sh/abhi-challenge4:2h"
    }

    stages {

        stage('Build image') {
            steps {
                sh '''
                    set -eu
                    docker build --platform linux/amd64 -t $IMAGE_NAME .
                '''
            }
        }

        stage('Push image') {
            steps {
                sh '''
                    set -eu
                    docker push $IMAGE_NAME
                '''
            }
        }

        stage('Deploy to Kubernetes') {
            steps {
                withCredentials([string(credentialsId: 'k8s-token', variable: 'KUBE_TOKEN')]) {
                    sh '''
                        set -eu

                        export KUBECONFIG=$WORKSPACE/kubeconfig

                        kubectl config set-cluster lab \
                          --server=https://kubernetes:6443 \
                          --insecure-skip-tls-verify=true

                        kubectl config set-credentials jenkins \
                          --token=$KUBE_TOKEN

                        kubectl config set-context lab \
                          --cluster=lab \
                          --user=jenkins

                        kubectl config use-context lab

                        echo "Applying Kubernetes manifests..."

                        kubectl apply -f k8s/deployment.yaml --validate=false
                        kubectl apply -f k8s/service.yaml --validate=false

                        echo "Waiting for rollout..."
                        kubectl rollout status deployment/myapp --timeout=120s || true
                    '''
                }
            }
        }
    }
}