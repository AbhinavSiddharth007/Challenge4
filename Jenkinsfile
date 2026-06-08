// pipeline {
//     agent any

//     // Credential ID is parameterized so the Kubernetes token can be rotated without
//     // editing the pipeline.
//     parameters {
//         string(
//             name: 'IMAGE_NAME',
//             defaultValue: 'ttl.sh/abhi-challenge4:2h',
//             description: 'Image tag pushed to ttl.sh (intentionally short-lived: the ":2h" tag expires 2h after push).'
//         )
//         string(
//             name: 'KUBE_TOKEN_CREDENTIAL_ID',
//             defaultValue: 'k8s-token',
//             description: 'ID of the Jenkins "Secret text" credential holding the Kubernetes ServiceAccount bearer token.'
//         )
//         string(
//             name: 'KUBE_API_SERVER',
//             defaultValue: 'https://kubernetes:6443',
//             description: 'Kubernetes API server URL.'
//         )
//         string(
//             name: 'KUBE_NAMESPACE',
//             defaultValue: 'default',
//             description: 'Namespace to deploy into.'
//         )
//         string(
//             name: 'POD_MANIFEST',
//             defaultValue: 'k8s/pod.yaml',
//             description: 'Path to the Pod manifest (adjust if yours lives at ./pod.yaml).'
//         )
//     }

//     environment {
//         IMAGE_NAME = "${params.IMAGE_NAME}"
//     }

//     stages {

//         stage('Build image') {
//             steps {
//                 sh '''
//                     set -eu
//                     docker build --platform linux/amd64 -t "$IMAGE_NAME" .
//                 '''
//             }
//         }

//         stage('Push image') {
//             steps {
//                 sh '''
//                     set -eu
//                     docker push "$IMAGE_NAME"
//                 '''
//             }
//         }

//         stage('Deploy to Kubernetes') {
//             steps {
//                 // Secret text credential -> bearer token. Jenkins masks $KUBE_TOKEN in the log.
//                 withCredentials([
//                     string(credentialsId: params.KUBE_TOKEN_CREDENTIAL_ID, variable: 'KUBE_TOKEN')
//                 ]) {
//                     // Single-quoted block + shell-level expansion => the token value is never
//                     // printed (only the literal $KUBE_TOKEN). Do NOT add `set -x` here.
//                     sh '''
//                         set -eu

//                         # Isolated kubeconfig for THIS build only, cleaned up on exit.
//                         export KUBECONFIG="$(pwd)/.kubeconfig.$$"
//                         trap 'rm -f "$KUBECONFIG"' EXIT
//                         : > "$KUBECONFIG"
//                         chmod 600 "$KUBECONFIG"

//                         # ---- Cluster entry (prefer in-cluster CA; else skip TLS verify) ----
//                         CA_FILE="/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
//                         if [ -f "$CA_FILE" ]; then
//                             echo "Using in-cluster CA for TLS verification."
//                             kubectl config set-cluster challenge-cluster \
//                                 --server="$KUBE_API_SERVER" \
//                                 --certificate-authority="$CA_FILE" \
//                                 --embed-certs=true
//                         else
//                             echo "WARNING: cluster CA not found; skipping TLS verification."
//                             kubectl config set-cluster challenge-cluster \
//                                 --server="$KUBE_API_SERVER" \
//                                 --insecure-skip-tls-verify=true
//                         fi

//                         kubectl config set-credentials jenkins-deployer --token="$KUBE_TOKEN"
//                         kubectl config set-context challenge \
//                             --cluster=challenge-cluster \
//                             --user=jenkins-deployer \
//                             --namespace="$KUBE_NAMESPACE"
//                         kubectl config use-context challenge

//                         # ---- Debug: context, cluster info, auth status --------------------
//                         echo "=== Current context ==="
//                         kubectl config current-context || true

//                         echo "=== Cluster info ==="
//                         kubectl cluster-info || true

//                         echo "=== Authentication status ==="
//                         kubectl auth whoami 2>/dev/null || echo "(kubectl auth whoami not supported on this version)"
//                         if kubectl auth can-i create pods -n "$KUBE_NAMESPACE" >/dev/null 2>&1; then
//                             echo "RBAC: ServiceAccount CAN create pods in $KUBE_NAMESPACE"
//                         else
//                             echo "RBAC: ServiceAccount CANNOT create pods in $KUBE_NAMESPACE (check its RoleBinding)"
//                         fi

//                         # ---- Deploy the Pod ------------------------------------------------
//                         # Recreate so the kubelet pulls the freshly-pushed image (the ttl.sh
//                         # tag is reused; imagePullPolicy: Always re-pulls on each start).
//                         echo "=== Applying Pod manifest ==="
//                         kubectl delete pod myapp -n "$KUBE_NAMESPACE" --ignore-not-found
//                         kubectl apply -n "$KUBE_NAMESPACE" -f "$POD_MANIFEST"

//                         # ---- Verify the Pod is Ready, with diagnostics on failure ---------
//                         echo "=== Waiting for Pod myapp to be Ready ==="
//                         kubectl wait --for=condition=Ready pod/myapp -n "$KUBE_NAMESPACE" --timeout=120s || {
//                             echo "Pod did not become Ready. Diagnostics:"
//                             kubectl describe pod myapp -n "$KUBE_NAMESPACE" || true
//                             echo "--- Recent events ---"
//                             kubectl get events -n "$KUBE_NAMESPACE" --sort-by=.lastTimestamp | tail -n 30 || true
//                             exit 1
//                         }

//                         echo "=== Pod is Ready ==="
//                         kubectl get pod myapp -n "$KUBE_NAMESPACE" -o wide
//                     '''
//                 }
//             }
//         }
//     }
// }

pipeline {
    agent any

    parameters {
        string(
            name: 'IMAGE_NAME',
            defaultValue: 'ttl.sh/abhi-challenge4:2h',
            description: 'Image tag pushed to ttl.sh (the ":2h" tag expires 2h after push).'
        )
        string(
            name: 'KUBE_TOKEN_CREDENTIAL_ID',
            defaultValue: 'k8s-token',
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
        string(
            name: 'POD_MANIFEST',
            defaultValue: 'pod.yaml',
            description: 'Path to the Pod manifest, relative to the repo root.'
        )
    }

    // Source EVERY value the shell needs from params here. params.* always carries its
    // default value -- even on a brand-new job's very first build -- whereas relying on
    // Jenkins to inject parameters as shell env vars fails on a first/recreated build
    // (which is exactly what bit the previous run: "KUBE_API_SERVER: parameter not set").
    environment {
        IMAGE_NAME      = "${params.IMAGE_NAME}"
        KUBE_API_SERVER = "${params.KUBE_API_SERVER}"
        KUBE_NAMESPACE  = "${params.KUBE_NAMESPACE}"
        POD_MANIFEST    = "${params.POD_MANIFEST}"
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
                    string(credentialsId: params.KUBE_TOKEN_CREDENTIAL_ID, variable: 'KUBE_TOKEN')
                ]) {
                    sh '''
                        set -eu

                        # Safety net: if any of these are ever empty, fall back to sane
                        # defaults so the script never dies on an unset variable again.
                        KUBE_API_SERVER="${KUBE_API_SERVER:-https://kubernetes:6443}"
                        KUBE_NAMESPACE="${KUBE_NAMESPACE:-default}"
                        POD_MANIFEST="${POD_MANIFEST:-pod.yaml}"

                        # Isolated kubeconfig for THIS build only, cleaned up on exit.
                        export KUBECONFIG="$(pwd)/.kubeconfig.$$"
                        trap 'rm -f "$KUBECONFIG"' EXIT
                        : > "$KUBECONFIG"
                        chmod 600 "$KUBECONFIG"

                        # ---- Cluster entry (prefer in-cluster CA; else skip TLS verify) ----
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

                        kubectl config set-credentials jenkins-deployer --token="$KUBE_TOKEN"
                        kubectl config set-context challenge \
                            --cluster=challenge-cluster \
                            --user=jenkins-deployer \
                            --namespace="$KUBE_NAMESPACE"
                        kubectl config use-context challenge

                        # ---- Debug: context, cluster info, auth status --------------------
                        echo "=== Current context ==="
                        kubectl config current-context || true

                        echo "=== Cluster info ==="
                        kubectl cluster-info || true

                        echo "=== Authentication status ==="
                        kubectl auth whoami 2>/dev/null || echo "(kubectl auth whoami not supported on this version)"
                        if kubectl auth can-i create pods -n "$KUBE_NAMESPACE" >/dev/null 2>&1; then
                            echo "RBAC: ServiceAccount CAN create pods in $KUBE_NAMESPACE"
                        else
                            echo "RBAC: ServiceAccount CANNOT create pods in $KUBE_NAMESPACE (check its RoleBinding)"
                        fi

                        # ---- Deploy the Pod ------------------------------------------------
                        echo "=== Applying Pod manifest ($POD_MANIFEST) ==="
                        kubectl delete pod myapp -n "$KUBE_NAMESPACE" --ignore-not-found
                        kubectl apply -n "$KUBE_NAMESPACE" -f "$POD_MANIFEST"

                        # ---- Verify the Pod is Ready, with diagnostics on failure ---------
                        echo "=== Waiting for Pod myapp to be Ready ==="
                        kubectl wait --for=condition=Ready pod/myapp -n "$KUBE_NAMESPACE" --timeout=120s || {
                            echo "Pod did not become Ready. Diagnostics:"
                            kubectl describe pod myapp -n "$KUBE_NAMESPACE" || true
                            echo "--- Recent events ---"
                            kubectl get events -n "$KUBE_NAMESPACE" --sort-by=.lastTimestamp | tail -n 30 || true
                            exit 1
                        }

                        echo "=== Pod is Ready ==="
                        kubectl get pod myapp -n "$KUBE_NAMESPACE" -o wide
                    '''
                }
            }
        }
    }
}