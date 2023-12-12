package builder

import (
	"encoding/json"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	dto_interface "github.com/Borislavv/video-streaming/internal/domain/dto/interface"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	"io"
	"net/http"
)

type AuthBuilder struct {
	logger logger_interface.Logger
}

func NewAuthBuilder(serviceContainer di_interface.ContainerManager) (*AuthBuilder, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	return &AuthBuilder{logger: loggerService}, nil
}

func (b *AuthBuilder) BuildAuthRequestDTOFromRequest(r *http.Request) (dto_interface.AuthRequest, error) {
	authDTO := &dto.AuthRequestDTO{}
	if err := json.NewDecoder(r.Body).Decode(authDTO); err != nil {
		if err == io.EOF {
			return nil, errors.NewRequestBodyIsEmptyError()
		}
		return nil, b.logger.LogPropagate(err)
	}
	return authDTO, nil
}
