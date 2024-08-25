# Stage 1: Build the application
FROM node:22-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy package.json and package-lock.json to the working directory
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application code to the working directory
COPY . .

# Build the application
RUN npm run build

# Stage 2: Serve the application using a lightweight server
FROM node:22-alpine

# Set the working directory
WORKDIR /app

# Copy the built application from the builder stage
COPY --from=builder /app ./

# Install the production dependencies only
RUN npm install --only=production

# Expose port 3000 to the outside world
EXPOSE 3000

# Start the Next.js application
CMD ["npm", "start"]
