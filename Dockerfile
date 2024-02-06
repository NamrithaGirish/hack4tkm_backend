# Use an official Golang runtime as a parent image
FROM golang:latest

# Set the working directory in the container
WORKDIR /go/src/app

# Copy the local package files to the container's workspace
COPY . .

# Install any needed packages
RUN go get -d -v ./...

# Install the Go application
RUN go install -v ./...

# Set environment variables
ENV PORT=8080

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the application
CMD ["hack4tkm"]  
