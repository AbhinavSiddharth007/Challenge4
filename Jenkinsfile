pipeline {
    agent any

    environment {
        IMAGE_NAME  = "ttl.sh/abhi-challenge4:2h"
        KUBE_SERVER = "https://kubernetes:6443"
    }

    stages {

        stage('Build image') {
            steps {
                sh '''
                    set -eu
                    docker build --platform linux/amd64 -t "$IMAGE_NAME" .
                '''
            }
        }

        stage('Push image') {
            steps {
                sh '''
                    set -eu
                    docker push "$IMAGE_NAME"
                '''
            }
        }

        stage('Deploy to Kubernetes') {
            steps {
                withCredentials([string(credentialsId: 'k8s-token', variable: 'KUBE_TOKEN')]) {
                    sh '''
                        set -eu

                        kubectl config set-cluster lab \
                          --server=$KUBE_SERVER \
                          --insecure-skip-tls-verify=true

                        kubectl config set-credentials jenkins \
                          --token=$KUBE_TOKEN

                        kubectl config set-context lab \
                          --cluster=lab \
                          --user=jenkins

                        kubectl config use-context lab

                        kubectl delete pod myapp --ignore-not-found
                        kubectl apply -f pod.yaml
                        kubectl wait --for=condition=Ready pod/myapp --timeout=120s
                    '''
                }
            }
        }
    }
}