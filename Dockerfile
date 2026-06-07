FROM golang:1.24 AS builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=builder /src/app .
CMD ["/app/app"]