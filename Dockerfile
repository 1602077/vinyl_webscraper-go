FROM golang:1.18 as builder
WORKDIR /app
COPY go/go.mod go/go.sum ./
RUN go mod download && go mod verify
COPY go/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o webscraper -a -installsuffix cgo ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/webscraper bin/
COPY .env input.txt .
COPY /sql ./sql
EXPOSE 8080
WORKDIR /app/bin
CMD ["./webscraper"]
