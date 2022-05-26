FROM golang:1.18 as builder
RUN mkdir /app
WORKDIR /app
COPY . .
RUN go mod download
RUN GOOS=linux CGO_ENABLED=0 go build -a -o /app/go-web-starter ./cmd/server

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/go-web-starter /app/go-web-starter
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/assets /app/assets
COPY --from=builder /app/sql /app/sql
ENTRYPOINT ["/app/go-web-starter"]
