# DEBUG

## Seeded bugs

### Bug 1: Readiness probe wrong port

Hypothesis:
Traffic was not routed because the readiness check never succeeded.

Verification:

```bash
kubectl describe pod
```

Fix:
Correct the `readinessProbe` port to `4444`.

### Bug 2: Missing resource limits

Hypothesis:
The pod could consume too much memory and make scheduling or eviction behavior unstable.

Verification:

```bash
kubectl describe pod
kubectl get pod -o yaml
```

Fix:
Add memory requests and limits to the container resources block.

### Bug 3: Service selector mismatch

Hypothesis:
The Service had no endpoints because its selector labels did not match the pod template labels.

Verification:

```bash
kubectl get endpoints
kubectl describe service myapp
```

Fix:
Make the Service selector match the Deployment labels exactly.

## Lessons learned

- Pods are the smallest runnable unit in Kubernetes, but a Deployment is usually the right abstraction when you want rollout control and replica management.
- Services give pods a stable virtual IP and DNS name, which is essential because pod IPs change whenever Kubernetes reschedules work.
- Readiness probes decide when traffic should be sent to a pod, while liveness probes decide when a container should be restarted.
- Resource requests and limits affect scheduling, stability, and autoscaling behavior. HPA also needs CPU requests to calculate utilization correctly.
