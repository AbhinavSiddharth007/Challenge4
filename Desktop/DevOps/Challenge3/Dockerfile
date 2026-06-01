# syntax=docker/dockerfile:1.7

FROM --platform=$BUILDPLATFORM golang:1.24 AS builder

WORKDIR /src

COPY go.mod ./
COPY main.go ./

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w" -o /out/main .

FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=builder /out/main /app/main

EXPOSE 4444

USER 65532:65532

CMD ["/app/main"]
