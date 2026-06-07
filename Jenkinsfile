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

                        export KUBECONFIG=$PWD/kubeconfig

                        cat > kubeconfig <<EOF
apiVersion: v1
kind: Config
clusters:
- cluster:
    insecure-skip-tls-verify: true
    server: $KUBE_SERVER
  name: lab
contexts:
- context:
    cluster: lab
    user: jenkins
  name: lab
current-context: lab
users:
- name: jenkins
  user:
    token: $KUBE_TOKEN
EOF

                        echo "Testing cluster connection..."
                        kubectl get nodes || true

                        echo "Applying deployment..."
                        kubectl apply -f k8s/deployment.yaml

                        kubectl rollout status deployment/myapp --timeout=120s || true
                    '''
                }
            }
        }
    }
}