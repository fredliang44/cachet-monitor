FROM golang:alpine as builder
LABEL maintainer="Fred Liang <info@fredliang.cn>"

# Set the Current Working Directory inside the container
WORKDIR /app
# Copy go mod and sum files
COPY go.mod go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
# Copy the source from the current directory to the Working Directory inside the container
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cachet-monitor

# Copy into a second stage container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /app/cachet-monitor /usr/local/bin/
RUN chmod +x /usr/local/bin/cachet-monitor
ENTRYPOINT ["/usr/local/bin/cachet-monitor", "-c/etc/cachet-monitor.config.json"]