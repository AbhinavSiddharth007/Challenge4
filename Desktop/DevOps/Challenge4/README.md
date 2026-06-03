# CS411 Challenge 4: Deploy Application to Kubernetes with Jenkins

This repository contains the Kubernetes deployment assets and Jenkins pipeline for the CS411 Challenge 4 submission.

The goal is to build the Go application image, push it to `ttl.sh`, and deploy the application into Kubernetes with a Jenkins pipeline.

## What’s Included

- `Jenkinsfile` - builds the image, pushes it to `ttl.sh`, authenticates to Kubernetes, and applies the manifests
- `k8s/deployment.yaml` - Kubernetes Deployment for `myapp`
- `k8s/service.yaml` - ClusterIP Service for in-cluster access
- `k8s/configmap.yaml` - Application settings injected into the pod
- `k8s/hpa.yaml` - HorizontalPodAutoscaler using `autoscaling/v2`
- `DEBUG.md` - seeded bug write-up and fixes
- `PROMPTS.md` - engineering decision log
- `.gitignore` - keeps build artifacts and local junk out of the repo

## Application Summary

- The application listens on port `4444`
- The root endpoint `/` is used by the readiness and liveness probes
- The Deployment consumes environment variables from a ConfigMap
- The Service exposes the app inside the cluster on port `4444`

## Jenkins Pipeline

The pipeline performs these steps:

1. Build the Docker image
2. Push the image to `ttl.sh`
3. Authenticate to the Kubernetes API at `https://kubernetes:6443`
4. Apply the ConfigMap, Deployment, Service, and HPA manifests
5. Wait for `kubectl rollout status deployment/myapp`

### Jenkins parameters

- `IMAGE_NAME` - fully qualified `ttl.sh` image reference
- `KUBE_TOKEN_CREDENTIAL_ID` - Jenkins Secret Text credential ID that stores the Kubernetes bearer token

## Kubernetes Resources

### Deployment

- Name: `myapp`
- Replicas: `1`
- Label selector: `app: myapp`
- Container port: `4444`
- Readiness probe: `GET /` on port `4444`
- Liveness probe: `GET /` on port `4444`
- Memory requests and limits are configured

### Service

- Type: `ClusterIP`
- Selector: `app: myapp`
- Port: `4444`
- Target port: `4444`

### ConfigMap

The ConfigMap is named `myapp-config` and provides example values such as:

- `APP_NAME=myapp`
- `PORT=4444`
- `LOG_LEVEL=info`

### HPA

- Targets the `myapp` Deployment
- Minimum replicas: `1`
- Maximum replicas: `3`
- CPU target: `70%`
- API version: `autoscaling/v2`

## Verification

After the pipeline runs, verify the deployment with:

```bash
kubectl get deployment
kubectl get pods
kubectl get svc
kubectl get hpa
kubectl get endpoints
kubectl rollout status deployment/myapp
```

Expected results:

- Deployment is available
- Pod is Ready
- Service has endpoints
- HPA is created successfully

## Notes

- The image is intended to be short-lived and hosted on `ttl.sh`
- The Kubernetes credential should be created in Jenkins as a Secret Text credential
- The pipeline assumes access to a Kubernetes API endpoint at `https://kubernetes:6443`
- The application source used for the Docker build lives in the repository’s existing challenge folder structure
