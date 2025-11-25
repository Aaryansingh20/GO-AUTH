# Use official Go image with your required version
FROM golang:1.25.4-alpine

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy all source code
COPY . .

# Build the application
RUN go build -o main .

# Expose the port your app runs on
EXPOSE 8080

# Run the application
CMD ["./main"]
