pipeline {
    agent any

    // The credential ID is parameterized (per the project's documented design) so the
    // Kubernetes token can be rotated/swapped without editing the pipeline.
    parameters {
        string(
            name: 'IMAGE_NAME',
            defaultValue: 'ttl.sh/abhi-challenge4:2h',
            description: 'Image tag pushed to ttl.sh (intentionally short-lived).'
        )
        string(
            name: 'KUBE_TOKEN_CREDENTIAL_ID',
            defaultValue: 'kube-token',
            description: 'ID of the Jenkins "Secret text" credential holding the Kubernetes ServiceAccount bearer token.'
        )
        string(
            name: 'KUBE_API_SERVER',
            defaultValue: 'https://kubernetes:6443',
            description: 'Kubernetes API server URL.'
        )
        string(
            name: 'KUBE_NAMESPACE',
            defaultValue: 'default',
            description: 'Namespace to deploy into.'
        )
    }

    environment {
        IMAGE_NAME = "${params.IMAGE_NAME}"
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
                // Secret text credential -> bearer token. Jenkins masks $KUBE_TOKEN in the log.
                withCredentials([
                    string(credentialsId: params.KUBE_TOKEN_CREDENTIAL_ID, variable: 'KUBE_TOKEN')
                ]) {
                    // NOTE: single-quoted heredoc + shell-level expansion means the token value
                    // is never printed in the Jenkins-echoed command (only the literal $KUBE_TOKEN).
                    // Do NOT add `set -x` here, and do NOT echo $KUBE_TOKEN.
                    sh '''
                        set -eu

                        # ---------------------------------------------------------------
                        # Build an isolated kubeconfig for THIS build only, then clean up.
                        # ---------------------------------------------------------------
                        export KUBECONFIG="$(pwd)/.kubeconfig.$$"
                        trap 'rm -f "$KUBECONFIG"' EXIT
                        : > "$KUBECONFIG"
                        chmod 600 "$KUBECONFIG"

                        # ---- Cluster entry --------------------------------------------
                        # Prefer the in-cluster CA when the agent runs inside Kubernetes.
                        # Otherwise skip TLS verification: acceptable for this short-lived
                        # ttl.sh challenge, but for production supply the cluster CA instead
                        # (see the "more secure variant" note in the accompanying explanation).
                        CA_FILE="/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
                        if [ -f "$CA_FILE" ]; then
                            echo "Using in-cluster CA for TLS verification."
                            kubectl config set-cluster challenge-cluster \
                                --server="$KUBE_API_SERVER" \
                                --certificate-authority="$CA_FILE" \
                                --embed-certs=true
                        else
                            echo "WARNING: cluster CA not found; skipping TLS verification."
                            kubectl config set-cluster challenge-cluster \
                                --server="$KUBE_API_SERVER" \
                                --insecure-skip-tls-verify=true
                        fi

                        # ---- Credentials: the bearer token from the Jenkins Secret text -
                        kubectl config set-credentials jenkins-deployer --token="$KUBE_TOKEN"

                        # ---- Context --------------------------------------------------
                        kubectl config set-context challenge \
                            --cluster=challenge-cluster \
                            --user=jenkins-deployer \
                            --namespace="$KUBE_NAMESPACE"
                        kubectl config use-context challenge

                        # ---------------------------------------------------------------
                        # Debugging: current context, cluster info, authentication status.
                        # All wrapped so every line prints even if one fails, giving a full
                        # picture before the (hard-gated) apply below.
                        # ---------------------------------------------------------------
                        echo "=== Current context ==="
                        kubectl config current-context || true

                        echo "=== Cluster info ==="
                        kubectl cluster-info || true

                        echo "=== Authentication status ==="
                        kubectl auth whoami 2>/dev/null || echo "(kubectl auth whoami not supported on this kubectl/server version)"
                        if kubectl auth can-i create deployments -n "$KUBE_NAMESPACE" >/dev/null 2>&1; then
                            echo "RBAC: ServiceAccount CAN create deployments in $KUBE_NAMESPACE"
                        else
                            echo "RBAC: ServiceAccount CANNOT create deployments in $KUBE_NAMESPACE (check its RoleBinding)"
                        fi

                        # ---- Apply manifests ------------------------------------------
                        echo "=== Applying manifests ==="
                        kubectl apply -n "$KUBE_NAMESPACE" -f k8s/deployment.yaml
                        # If you also keep service/configmap/hpa manifests, apply them here, e.g.:
                        #   kubectl apply -n "$KUBE_NAMESPACE" -f k8s/service.yaml -f k8s/configmap.yaml -f k8s/hpa.yaml

                        # ---- Verify rollout -------------------------------------------
                        echo "=== Waiting for rollout ==="
                        kubectl rollout status deployment/myapp -n "$KUBE_NAMESPACE" --timeout=120s
                    '''
                }
            }
        }
    }
}
