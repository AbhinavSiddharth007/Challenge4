# PROMPTS

## What I asked for

I asked for a production-style Docker setup around the existing Go service: build the binary inside a multi-stage image, publish it to `ttl.sh`, and then deploy that image onto the remote Docker VM without changing the required port `4444`.

## What I decided

I kept the runtime small by using a distroless final image and a static Go build. I also forced the Jenkins build to target `linux/amd64`, because the remote machine that runs the challenge verification is x86_64 and I didn’t want a local laptop architecture to leak into the release artifact.

## What I pushed back on

I avoided a single-stage Dockerfile because it would have bundled build tooling into the final image. I also rejected any deploy flow that assumed the Docker VM was reachable locally from Jenkins, since the real handoff happens over SSH and the image needs to be pulled on the target host before the container starts.
