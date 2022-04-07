# golang builder img
FROM golang:1.18 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo ./cmd/main.go

# prod img
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/data ./data
COPY --from=builder /app/.env .
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
