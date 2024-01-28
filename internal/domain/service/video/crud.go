package video

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/builder/interface"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	dtointerface "github.com/Borislavv/video-streaming/internal/domain/dto/interface"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	repositoryinterface "github.com/Borislavv/video-streaming/internal/domain/repository/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	resourceinterface "github.com/Borislavv/video-streaming/internal/domain/service/resource/interface"
	validatorinterface "github.com/Borislavv/video-streaming/internal/domain/validator/interface"
)

type CRUDService struct {
	ctx             context.Context
	logger          loggerinterface.Logger
	builder         builderinterface.Video
	validator       validatorinterface.Video
	repository      repositoryinterface.Video
	resourceService resourceinterface.CRUD
}

func NewCRUDService(serviceContainer diinterface.ContainerManager) (*CRUDService, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	ctx, err := serviceContainer.GetCtx()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	videoBuilder, err := serviceContainer.GetVideoBuilder()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	videoValidator, err := serviceContainer.GetVideoValidator()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	videoRepository, err := serviceContainer.GetVideoRepository()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	resourceCRUDService, err := serviceContainer.GetResourceCRUDService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return &CRUDService{
		ctx:             ctx,
		logger:          loggerService,
		builder:         videoBuilder,
		validator:       videoValidator,
		repository:      videoRepository,
		resourceService: resourceCRUDService,
	}, nil
}

// Get - will fetch a single video aggregate by ID and specified user.
// Access check to video is unnecessary because the query will fetch video only for specified user.
func (s *CRUDService) Get(req dtointerface.GetVideoRequest) (*agg.Video, error) {
	// validation of input request
	if err := s.validator.ValidateGetRequestDTO(req); err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	// fetching a video by id and user
	video, err := s.repository.FindOneByID(s.ctx, req)
	if err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	return video, nil
}

// List - will fetch a video list of aggregates by given request and specified user.
// Access check to video is unnecessary because the query will fetch a video list only for specified user.
func (s *CRUDService) List(req dtointerface.ListVideoRequest) (list []*agg.Video, total int64, err error) {
	// validation of input request
	if err = s.validator.ValidateListRequestDTO(req); err != nil {
		return nil, 0, s.logger.LogPropagate(err)
	}

	// fetching a video list by request params. and user
	list, total, err = s.repository.FindList(s.ctx, req)
	if err != nil {
		return nil, 0, s.logger.LogPropagate(err)
	}

	return list, total, err
}

// Create - will make a new video by given request for specified user. Have an access check for resource
// which exists into the request.
func (s *CRUDService) Create(req dtointerface.CreateVideoRequest) (*agg.Video, error) {
	// validation of input request
	if err := s.validator.ValidateCreateRequestDTO(req); err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	// building an aggregate
	videoAgg, err := s.builder.BuildAggFromCreateRequestDTO(req)
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

// Update - will change the video by given request. Have an access check for video.resource.
func (s *CRUDService) Update(req dtointerface.UpdateVideoRequest) (*agg.Video, error) {
	// validation of input request
	if err := s.validator.ValidateUpdateRequestDTO(req); err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	// building an aggregate
	videoAgg, err := s.builder.BuildAggFromUpdateRequestDTO(req)
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

// Delete - will remove the video from the storage.
func (s *CRUDService) Delete(req dtointerface.DeleteVideoRequest) (err error) {
	// validation of input request
	if err = s.validator.ValidateDeleteRequestDTO(req); err != nil {
		return s.logger.LogPropagate(err)
	}

	// fetching a video which will be deleted
	videoAgg, err := s.repository.FindOneByID(s.ctx, req)
	if err != nil {
		return s.logger.LogPropagate(err)
	}

	// the resource must be removing first
	q := dto.NewResourceDeleteRequestDTO(videoAgg.Resource.ID, req.GetUserID())
	if err = s.resourceService.Delete(q); err != nil {
		return s.logger.LogPropagate(err)
	}

	// video removing
	if err = s.repository.Remove(s.ctx, videoAgg); err != nil {
		return s.logger.LogPropagate(err)
	}

	return nil
}
