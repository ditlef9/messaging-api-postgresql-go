# Dockerfile
FROM golang:1.23.2

WORKDIR /app
COPY . .

# Set environment variables for DB configuration
ENV DB_HOST=host.docker.internal
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=root
ENV DB_NAME=msg-dev

# Build
RUN go build -o main .

EXPOSE 8080

# Run
CMD ["./main"]