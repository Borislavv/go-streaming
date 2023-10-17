package video

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
	"github.com/Borislavv/video-streaming/internal/domain/service/resource"
	"github.com/Borislavv/video-streaming/internal/domain/validator"
)

type CRUDService struct {
	ctx             context.Context
	logger          logger.Logger
	builder         builder.Video
	validator       validator.Video
	repository      repository.Video
	resourceService resource.CRUD
}

func NewCRUDService(
	ctx context.Context,
	logger logger.Logger,
	builder builder.Video,
	validator validator.Video,
	repository repository.Video,
	resourceService resource.CRUD,
) *CRUDService {
	return &CRUDService{
		ctx:             ctx,
		logger:          logger,
		builder:         builder,
		validator:       validator,
		repository:      repository,
		resourceService: resourceService,
	}
}

func (s *CRUDService) Get(reqDTO dto.GetVideoRequest) (*agg.Video, error) {
	// validation of input request
	if err := s.validator.ValidateGetRequestDTO(reqDTO); err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	video, err := s.repository.Find(s.ctx, reqDTO.GetID())
	if err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	return video, nil
}

func (s *CRUDService) List(reqDTO dto.ListVideoRequest) (list []*agg.Video, total int64, err error) {
	// validation of input request
	if err = s.validator.ValidateListRequestDTO(reqDTO); err != nil {
		return nil, 0, s.logger.LogPropagate(err)
	}

	list, total, err = s.repository.FindList(s.ctx, reqDTO)
	if err != nil {
		return nil, 0, s.logger.LogPropagate(err)
	}

	return list, total, err
}

func (s *CRUDService) Create(reqDTO dto.CreateVideoRequest) (*agg.Video, error) {
	// validation of input request
	if err := s.validator.ValidateCreateRequestDTO(reqDTO); err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	// building an aggregate
	videoAgg, err := s.builder.BuildAggFromCreateRequestDTO(reqDTO)
	if err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	// validation of an aggregate
	if err = s.validator.ValidateAggregate(videoAgg); err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	// saving an aggregate into storage
	videoAgg, err = s.repository.Insert(s.ctx, videoAgg)
	if err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	return videoAgg, nil
}

func (s *CRUDService) Update(reqDTO dto.UpdateVideoRequest) (*agg.Video, error) {
	// validation of input request
	if err := s.validator.ValidateUpdateRequestDTO(reqDTO); err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	// building an aggregate
	videoAgg, err := s.builder.BuildAggFromUpdateRequestDTO(reqDTO)
	if err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	// validation of an aggregate
	if err = s.validator.ValidateAggregate(videoAgg); err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	// saving updated aggregate into storage
	videoAgg, err = s.repository.Update(s.ctx, videoAgg)
	if err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	return videoAgg, nil
}

func (s *CRUDService) Delete(reqDTO dto.DeleteVideoRequest) (err error) {
	// validation of input request
	if err = s.validator.ValidateDeleteRequestDTO(reqDTO); err != nil {
		return s.logger.LogPropagate(err)
	}

	// fetching a video which will be deleted
	videoAgg, err := s.repository.Find(s.ctx, reqDTO.GetID())
	if err != nil {
		return s.logger.LogPropagate(err)
	}

	// resource removing first
	if err = s.resourceService.Delete(&dto.DeleteResourceRequestDTO{ID: videoAgg.Resource.ID}); err != nil {
		return s.logger.LogPropagate(err)
	}

	// video removing
	if err = s.repository.Remove(s.ctx, videoAgg); err != nil {
		return s.logger.LogPropagate(err)
	}

	return nil
}
