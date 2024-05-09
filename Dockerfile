# Start from the latest golang base image
FROM golang:1.22.0-alpine AS build-env

# Add Maintainer Info
LABEL maintainer="DeanXu2357 <dean.xu.2357@gmail.com>"

RUN apk --no-cache add curl

# Install delve
RUN go install github.com/go-delve/delve/cmd/dlv@latest

FROM build-env

# Set the Current Working Directory inside the Docker container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the Docker container
COPY . .

# Command to run the executable
CMD ["go", "run", "main.go", "serve"]
