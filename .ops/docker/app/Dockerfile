FROM golang:1.18

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Install FFmpeg including ffprobe
RUN apt-get update && apt-get install -y ffmpeg

# Set the working directory inside the container
WORKDIR /streaming_root

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download and install Go dependencies
RUN go mod download

# Copy the rest of the project to the working directory
COPY . .