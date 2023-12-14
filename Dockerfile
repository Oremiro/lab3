# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Build the Golang application
RUN go build -o file-storage-app .

# Expose port 8080 for the application
EXPOSE 8080

# Command to run the executable
CMD ["./file-storage-app"]
