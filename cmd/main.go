package main

import (
	"github.com/Borislavv/video-streaming/internal/app/service/stream"
	"github.com/Borislavv/video-streaming/internal/app/service/video"
	"github.com/Borislavv/video-streaming/internal/infrastructure/server/socket"
	"log"
)

func main() {
	manager := video.NewVideoManagerService()
	streamer := stream.NewStreamingService(manager)
	server := socket.NewSocketServer(streamer)

	if err := server.Listen(); err != nil {
		log.Fatalln(err)
	}
}
