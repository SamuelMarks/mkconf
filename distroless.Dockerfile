FROM golang:1.22-bookworm AS builder
WORKDIR /app
COPY . .
RUN go build -o app
FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=builder /app/app /app/app
ENTRYPOINT ["/app/app"]
