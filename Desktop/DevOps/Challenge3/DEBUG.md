# DEBUG — Architecture Mismatch

## Failure hypotheses

1. The image was built for the wrong CPU family, most likely `arm64`, while the Docker VM expects `amd64`.
2. The Go binary inside the image was produced without pinning the target platform, so the build host quietly shaped the artifact.
3. The runtime image was fine, but the executable copied into it could not be launched on the remote machine.

## How I checked

```bash
docker image inspect ttl.sh/abhi-challenge3:2h --format '{{.Architecture}}'
file main
docker inspect ttl.sh/abhi-challenge3:2h
```

To prove the fix locally, I would rebuild with the target platform pinned and then inspect the image again:

```bash
docker build --platform linux/amd64 -t ttl.sh/abhi-challenge3:2h .
docker image inspect ttl.sh/abhi-challenge3:2h --format '{{.Architecture}}'
```

## Fix

I made the Docker build explicitly target `linux/amd64` and cross-compiled the Go binary with `CGO_ENABLED=0`, `GOOS=linux`, and `GOARCH=amd64`. That keeps the release image aligned with the x86_64 Docker VM instead of whatever architecture happened to be available on the machine that ran Jenkins.

## Lesson

A successful build only matters if the artifact can actually execute on the host you deploy to.
