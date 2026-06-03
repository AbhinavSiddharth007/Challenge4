# Reflection — Docker build and deploy

## What did I do?

I took the Go service and turned it into a containerized release pipeline: the application now builds inside a multi-stage Dockerfile, the image gets pushed to `ttl.sh`, and Jenkins uses SSH to start the container on the remote Docker VM on port `4444`.

## What was the most interesting part?

The subtle part was making the deployment deterministic across machines. The same source code can still fail at runtime if the image is built for the wrong architecture, so I had to treat `linux/amd64` as part of the release contract rather than an implementation detail.

## What felt tricky?

The Jenkins deploy step was the easiest place to make a mistake because it crosses three shells at once: Jenkins, the local shell in the job, and the remote shell on the Docker VM. A small quoting mistake there can look like a container bug even when the real issue is just command interpolation.

## What would I keep the same next time?

I would still start with a static Go binary, a slim runtime image, and an explicit target platform. That combination keeps the image small, the runtime predictable, and the deployment path close to what the verifier expects.
