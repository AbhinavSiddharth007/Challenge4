# PROMPTS

## 1. What I asked the agent

- Build a Dockerfile for a Go service that listens on port 4444.
- Make the image production-ready with a multi-stage build and a static binary.
- Create a Jenkins pipeline that builds, pushes to `ttl.sh`, and deploys to a remote Docker VM.
- Add the required documentation files for decisions and debugging notes.

## 2. What decisions were made

- Chose a multi-stage Dockerfile so the final image only contains the runtime binary.
- Used `CGO_ENABLED=0` and cross-compilation flags so the binary is static and portable.
- Used `gcr.io/distroless/static-debian12:nonroot` as the final image to keep the runtime small.
- Kept the service on port 4444 end-to-end to match the challenge requirements.
- Built the image as `linux/amd64` in Jenkins so it can run on an x86_64 Docker VM.
- Used `ttl.sh` because it does not require auth for short-lived challenge images.

## 3. What I pushed back on / corrected

- Rejected a single-stage Dockerfile because it would ship build tools in the final image.
- Corrected the build target so the deployed image is amd64 instead of an ARM-only image.
- Kept the deploy stage as an SSH-based remote `docker pull` + `docker run` flow instead of assuming local access to the Docker VM.
- Kept the port mapping at `4444:4444` to avoid mismatches between the container and the verifier.
