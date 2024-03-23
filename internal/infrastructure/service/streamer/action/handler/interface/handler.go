package handlerinterface

import (
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/model"
	"sync"
)

type ActionsHandler interface {
	Handle(wg *sync.WaitGroup, actionsCh <-chan model.Action)
}
