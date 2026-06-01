# DEBUG

## Scenario

`exec format error` on the Docker VM, which is `x86_64`.

## 1. Hypotheses, ranked

1. The image was built for the wrong CPU architecture, most likely `arm64` from an Apple Silicon build host.
2. The final image was produced from the host architecture instead of the target architecture, so the container binary does not match the VM.
3. The binary was dynamically linked or otherwise built in a way that made it incompatible with the runtime image.

## 2. Verification commands

```bash
docker image inspect ttl.sh/abhi-challenge3:2h --format '{{.Architecture}}'
file main
docker inspect ttl.sh/abhi-challenge3:2h
```

If I need to confirm what the build produced locally:

```bash
docker build --platform linux/amd64 -t ttl.sh/abhi-challenge3:2h .
docker image inspect ttl.sh/abhi-challenge3:2h --format '{{.Architecture}}'
```

## 3. Fix

- Build the image for the Docker VM explicitly:

```bash
docker build --platform linux/amd64 -t ttl.sh/abhi-challenge3:2h .
```

- Cross-compile the Go binary inside Docker:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main
```

- Keep the runtime stage aligned with the target architecture.

## 4. Lesson

A built image does not guarantee runtime compatibility with the host architecture.
