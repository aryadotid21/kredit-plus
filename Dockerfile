# Use the official Golang 1.19 image as the base image
FROM golang:1.20

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules and download the dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the project files into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Set the environment variables for MongoDB
COPY .env .

# Start the application
CMD ["./main"]