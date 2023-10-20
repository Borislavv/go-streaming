package builder

import (
	"encoding/json"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"io"
	"net/http"
)

type AuthBuilder struct {
	logger logger.Logger
}

func NewAuthBuilder(logger logger.Logger) *AuthBuilder {
	return &AuthBuilder{logger: logger}
}

func (b *AuthBuilder) BuildAuthRequestDTOFromRequest(r *http.Request) (dto.AuthRequest, error) {
	authDTO := &dto.AuthRequestDTO{}
	if err := json.NewDecoder(r.Body).Decode(authDTO); err != nil {
		if err == io.EOF {
			return nil, errors.NewRequestBodyIsEmptyError()
		}
		return nil, b.logger.LogPropagate(err)
	}
	return authDTO, nil
}
