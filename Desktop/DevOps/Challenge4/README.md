# CS411 Challenge 4

This folder holds the deployment side of the challenge. I kept the actual application code where it already lived and focused this part of the repository on the Kubernetes and Jenkins pieces that turn that app into something deployable.

The workflow is simple on purpose: Jenkins builds the image, pushes it to `ttl.sh`, talks to the Kubernetes API with a Secret Text credential, and then applies the manifests that make up the release.

## How it fits together

The app listens on port `4444`, so every Kubernetes object is wired around that same port. The Deployment runs one replica, injects config from a ConfigMap, and uses readiness and liveness probes on `/` so Kubernetes only sends traffic when the app is actually responding. The Service gives the pod a stable in-cluster address, and the HPA keeps the replica count between `1` and `3` based on CPU usage.

## Files worth knowing

- `Jenkinsfile` drives build, push, apply, and rollout verification
- `k8s/deployment.yaml` defines the `myapp` Deployment
- `k8s/service.yaml` exposes the app on port `4444`
- `k8s/configmap.yaml` stores basic runtime settings
- `k8s/hpa.yaml` enables autoscaling with `autoscaling/v2`
- `DEBUG.md` explains the seeded Kubernetes bugs and their fixes
- `PROMPTS.md` records the decisions I made while solving the challenge

## Jenkins details

The pipeline expects two main inputs:

- `IMAGE_NAME`, which defaults to a short-lived `ttl.sh` tag
- `KUBE_TOKEN_CREDENTIAL_ID`, which points at the Jenkins Secret Text credential for the Kubernetes bearer token

It connects to `https://kubernetes:6443`, applies the manifests with `kubectl apply -f`, and finishes by waiting on `kubectl rollout status deployment/myapp`.

## What to verify

After the pipeline runs, the basic checks should all line up:

```bash
kubectl get deployment
kubectl get pods
kubectl get svc
kubectl get hpa
kubectl get endpoints
```

The deployment should be available, the pod should be Ready, the Service should have endpoints, and the HPA should exist with the expected replica bounds.

## Small but important details

- The Service selector must match `app: myapp`
- The probes must use port `4444`
- The Deployment includes memory requests and limits
- The ConfigMap is mounted through environment variables rather than hardcoded values
- The image is meant to be temporary because `ttl.sh` is used for the registry
