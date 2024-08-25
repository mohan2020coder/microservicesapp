# # Base Go image for building the application
# FROM golang:1.22-alpine as builder

# # Create the application directory
# RUN mkdir /app

# # Copy all the project files to the container
# COPY . /app

# # Set the working directory
# WORKDIR /app

# # Build the Go application binary
# RUN CGO_ENABLED=0 go build -o authApp ./cmd/api

# # Set permissions to make the binary executable
# RUN chmod +x /app/authApp

# # Final image: a minimal Alpine image to run the app
# FROM alpine:latest

# # Create the application directory in the final image
# RUN mkdir /app

# # Copy the built binary from the builder stage
# COPY --from=builder /app/authApp /app

# # Copy the wait script (e.g., wait-for-it.sh) to the final image
# COPY wait-for-it.sh /usr/local/bin/

# # Set permissions to make the wait script executable
# RUN chmod +x /usr/local/bin/wait-for-it.sh

# # Command to start the application, ensuring Postgres is ready
# CMD ["sh", "-c", "/usr/local/bin/wait-for-it.sh postgres:5432 -- /app/authApp"]


FROM alpine:latest

RUN mkdir /app

COPY authApp /app

CMD [ "/app/authApp"]