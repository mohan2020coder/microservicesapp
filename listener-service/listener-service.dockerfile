# # Build stage
# FROM golang:1.22-alpine as builder

# WORKDIR /app

# COPY . .

# # Adjust the build command to specify the main package
# RUN CGO_ENABLED=0 go build -o listenerApp ./cmd/api

# RUN chmod +x /app/listenerApp

# # Final stage
# FROM alpine:3.18

# RUN mkdir /app


# # Copy the built binary from the builder stage
# COPY --from=builder /app/listenerApp /app

# # Copy wait-for-it.sh script to the final image
# COPY wait-for-it.sh /usr/local/bin/

# # Ensure the binary and script are executable
# RUN chmod +x /app/listenerApp /usr/local/bin/wait-for-it.sh

# # Command to wait for RabbitMQ and then start the listenerApp
# CMD ["sh", "-c", "/usr/local/bin/wait-for-it.sh rabbitmq:5672 -- /app/listenerApp"]
# FROM golang:1.22-alpine as builder

# WORKDIR /app

# COPY . .

# RUN CGO_ENABLED=0 go build -o /app/listenerApp ./cmd/api

# RUN chmod +x /app/listenerApp

# FROM alpine:latest

# WORKDIR /app

# COPY --from=builder /app/listenerApp /app/listenerApp

# RUN chmod +x /app/listenerApp

# ENTRYPOINT ["/app/listenerApp"]
FROM alpine:latest

RUN mkdir /app

COPY listenerApp /app

CMD [ "/app/listenerApp"]