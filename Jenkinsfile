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
                withCredentials([
                    file(credentialsId: 'kubeconfig-file', variable: 'KUBECONFIG_FILE')
                ]) {
                    sh '''
                        set -eu

                        export KUBECONFIG=$KUBECONFIG_FILE

                        echo "Checking cluster connection..."
                        kubectl get nodes

                        echo "Applying Kubernetes manifests..."
                        kubectl apply -f k8s/deployment.yaml

                        echo "Restarting pod (optional)..."
                        kubectl delete pod myapp --ignore-not-found
                        kubectl apply -f k8s/pod.yaml

                        echo "Waiting for pod to be ready..."
                        kubectl wait --for=condition=Ready pod/myapp --timeout=120s
                    '''
                }
            }
        }
    }
}