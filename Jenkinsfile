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
    server: https://kubernetes:6443
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
    token: ${KUBE_TOKEN}
EOF

                kubectl get nodes || true

                kubectl apply -f k8s/deployment.yaml
                kubectl rollout status deployment/myapp --timeout=120s || true
            '''
        }
    }
}