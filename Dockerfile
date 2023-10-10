FROM golang:1.20

# Set the working directory
WORKDIR /app

# Copy the source code to the working directory
COPY . .

# Install the dependencies
RUN go mod tidy

# Build the application
RUN go build -o main .

# Expose the port that the application listens on
EXPOSE 80

# Start the application
CMD ["go","run","main.go"]