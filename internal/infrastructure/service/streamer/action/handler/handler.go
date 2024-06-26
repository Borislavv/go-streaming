package handler

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	strategyinterface "github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/handler/strategy/interface"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/model"
	"sync"
)

type WebSocketActionsHandler struct {
	ctx              context.Context
	logger           loggerinterface.Logger
	actionStrategies []strategyinterface.ActionStrategy
}

func NewWebSocketActionsHandler(serviceContainer diinterface.ServiceContainer) (*WebSocketActionsHandler, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	ctx, err := serviceContainer.GetCtx()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	strategies, err := serviceContainer.GetWebSocketHandlerStrategies()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return &WebSocketActionsHandler{
		ctx:              ctx,
		logger:           loggerService,
		actionStrategies: strategies,
	}, nil
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
