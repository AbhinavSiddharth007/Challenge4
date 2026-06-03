# PROMPTS

## GOAL

Deploy the existing Go application to Kubernetes with Jenkins, using `ttl.sh` for the container image and `kubectl apply -f` for the manifests.

## QUESTIONS ASKED

- Should I use a Pod or a Deployment?
- Why use a Service instead of exposing the pod directly?
- What is the difference between readiness and liveness probes?
- Why configure resource requests and limits?
- Why use a ConfigMap for application settings?
- How should Jenkins authenticate to the Kubernetes API?

## DECISIONS MADE

- Chose a Deployment instead of a bare Pod so the application can be rolled out and rescheduled safely.
- Added a ClusterIP Service so the app has a stable in-cluster endpoint.
- Added readiness and liveness probes on `/` port `4444` to match the Go service.
- Added a ConfigMap and injected it into the pod so runtime settings can be managed separately from the image.
- Added memory requests and limits, plus a CPU request, so the scheduler and HPA have the data they need.
- Used `autoscaling/v2` for the HPA with a 70 percent CPU target.
- Kept `ttl.sh` as the registry because it fits the short-lived challenge workflow and does not need registry auth.
- Parameterized the Jenkins credential ID so the Kubernetes token can be swapped without editing the pipeline.

## PUSHBACK / CORRECTIONS

- The first pass treated the workload like a single pod, but that would not satisfy the grading rubric or provide a clean rollout story, so it was updated to a Deployment.
- The initial manifest omitted resource requests, which would have made scheduling and autoscaling less reliable, so those were added back.
- A mismatched Service selector would have left the Service with no endpoints, so the labels were aligned exactly with the Deployment template.
- The first probe pass used the wrong port, which would have kept the pod unready even though the container was healthy, so the probes were corrected to `4444`.
