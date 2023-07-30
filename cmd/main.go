package main

import (
	"github.com/Borislavv/video-streaming/internal/app/service/stream"
	"github.com/Borislavv/video-streaming/internal/app/service/stream/read"
	"github.com/Borislavv/video-streaming/internal/infrastructure/server/socket"
	"log"
)

func main() {
	reader := read.NewReadingService()
	streamer := stream.NewStreamingService(reader)
	server := socket.NewSocketServer(streamer)

	if err := server.Listen(); err != nil {
		log.Fatalln(err)
	}
}
