package builder

import (
	"encoding/json"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	dtointerface "github.com/Borislavv/video-streaming/internal/domain/dto/interface"
	"github.com/Borislavv/video-streaming/internal/domain/errtype"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	"io"
	"net/http"
)

type AuthBuilder struct {
	logger loggerinterface.Logger
}

func NewAuthBuilder(serviceContainer diinterface.ContainerManager) (*AuthBuilder, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	return &AuthBuilder{logger: loggerService}, nil
}

func (b *AuthBuilder) BuildAuthRequestDTOFromRequest(r *http.Request) (dtointerface.AuthRequest, error) {
	authDTO := &dto.AuthRequestDTO{}
	if err := json.NewDecoder(r.Body).Decode(authDTO); err != nil {
		if err == io.EOF {
			return nil, errtype.NewRequestBodyIsEmptyError()
		}
		return nil, b.logger.LogPropagate(err)
	}
	return authDTO, nil
}
