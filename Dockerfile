FROM golang:latest AS builder
WORKDIR /app
COPY . .
# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w" -a -installsuffix cgo -o main .

FROM alpine:latest AS task
WORKDIR /root/
# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app .
CMD ["./main"]