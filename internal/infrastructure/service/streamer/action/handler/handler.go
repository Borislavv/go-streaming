package handler

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/handler/strategy"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/model"
	"sync"
)

type WebSocketActionsHandler struct {
	ctx              context.Context
	logger           logger.Logger
	actionStrategies []strategy.ActionStrategy
}

func NewWebSocketActionsHandler(
	ctx context.Context,
	logger logger.Logger,
	actionStrategies []strategy.ActionStrategy,
) *WebSocketActionsHandler {
	return &WebSocketActionsHandler{
		ctx:              ctx,
		logger:           logger,
		actionStrategies: actionStrategies,
	}
}

func (h *WebSocketActionsHandler) Handle(wg *sync.WaitGroup, actionsCh <-chan model.Action) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for action := range actionsCh {
			for _, actionStrategy := range h.actionStrategies {
				if actionStrategy.IsAppropriate(action) {
					if err := actionStrategy.Do(action); err != nil {
						h.logger.Error(err)
						break
					}
				}
			}
		}
	}()
}
