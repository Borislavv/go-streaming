package mongodb

import (
	"context"
	"errors"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/errs"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"sync"
	"time"
)

const VideosCollection = "videos"

var (
	VideoNotFoundByIdError       = errs.NewNotFoundError("video", "id")
	VideoNotFoundByResourceError = errs.NewNotFoundError("video", "resource")
)

type VideoRepository struct {
	db      *mongo.Collection
	mu      *sync.Mutex
	logger  logger.Logger
	timeout time.Duration
}

func NewVideoRepository(db *mongo.Database, logger logger.Logger, timeout time.Duration) *VideoRepository {
	return &VideoRepository{
		db:      db.Collection(VideosCollection),
		logger:  logger,
		mu:      &sync.Mutex{},
		timeout: timeout,
	}
}

func (r *VideoRepository) Find(ctx context.Context, id vo.ID) (*agg.Video, error) {
	qCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	video := &agg.Video{}
	if err := r.db.FindOne(qCtx, bson.M{"_id": bson.M{"$eq": id.Value}}).Decode(video); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, r.logger.InfoPropagate(VideoNotFoundByIdError)
		}
		return nil, r.logger.ErrorPropagate(err)
	}

	return video, nil
}

func (r *VideoRepository) FindByResource(ctx context.Context, resource *agg.Resource) (*agg.Video, error) {
	qCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	video := &agg.Video{}
	if err := r.db.FindOne(qCtx, bson.M{"resource": bson.M{"$eq": resource.Resource}}).Decode(video); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, VideoNotFoundByResourceError
		}
		return nil, r.logger.ErrorPropagate(err)
	}

	return video, nil
}

func (r *VideoRepository) FindList(ctx context.Context, dto dto.ListRequest) ([]*agg.Video, error) {
	qCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	filter := bson.M{}

	if dto.GetName() != "" {
		filter["name"] = primitive.Regex{Pattern: dto.GetName(), Options: "i"}
	}
	if !dto.GetCreatedAt().IsZero() {
		y := dto.GetCreatedAt().Year()
		m := dto.GetCreatedAt().Month()
		d := dto.GetCreatedAt().Day()

		filter["createdAt"] = bson.M{
			"$gt":  time.Date(y, m, d, 0, 0, 0, 0, time.UTC),
			"$lte": time.Date(y, m, d, 23, 59, 59, 0, time.UTC),
		}
	} else if !dto.GetFrom().IsZero() || !dto.GetTo().IsZero() {
		createdAtFilter := bson.M{}
		if !dto.GetFrom().IsZero() {
			createdAtFilter["$gt"] = dto.GetFrom()
		}
		if !dto.GetTo().IsZero() {
			createdAtFilter["$lte"] = dto.GetTo()
		}
		filter["createdAt"] = createdAtFilter
	}

	opts := options.Find().
		SetSkip((int64(dto.GetPage()) - 1) * int64(dto.GetLimit())).
		SetLimit(int64(dto.GetLimit()))

	videos := []*agg.Video{}
	cursor, err := r.db.Find(qCtx, filter, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return videos, nil
		}
		return nil, r.logger.ErrorPropagate(err)
	}
	defer func() { _ = cursor.Close(qCtx) }()

	if err = cursor.All(qCtx, &videos); err != nil {
		return nil, r.logger.ErrorPropagate(err)
	}

	return videos, nil
}

func (r *VideoRepository) Insert(ctx context.Context, video *agg.Video) (*agg.Video, error) {
	qCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	res, err := r.db.InsertOne(qCtx, video, options.InsertOne())
	if err != nil {
		return nil, r.logger.ErrorPropagate(err)
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		return r.Find(qCtx, vo.ID{Value: oid})
	}

	return nil, r.logger.CriticalPropagate("unable to store 'video' or retrieve inserted 'id'")
}

func (r *VideoRepository) Update(ctx context.Context, video *agg.Video) (*agg.Video, error) {
	qCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	res, err := r.db.UpdateByID(qCtx, video.ID.Value, bson.M{"$set": video})
	if err != nil {
		return nil, r.logger.ErrorPropagate(err)
	}

	// check the record is really updated
	if res.ModifiedCount > 0 {
		return r.Find(qCtx, video.ID)
	}

	// if changes is not exists, then return the original data
	return video, nil
}

func (r *VideoRepository) Has(ctx context.Context, video *agg.Video) (bool, error) {
	qCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	// TODO Stopped here, validation must be improved, probably need to migrate to store resource id instead of entity.Resource
	filter := bson.M{}
	if video.Name != "" && !reflect.DeepEqual(video.Resource, entity.Resource{}) {
		filter["$or"] = []bson.M{
			{"name": video.Name},
			{"resource": video.Resource},
		}
	} else if video.Name != "" {
		filter["name"] = video.Name
	} else if reflect.DeepEqual(video.Resource, entity.Resource{}) {
		filter["resource"] = video.Resource
	}

	foundVideo := &agg.Video{}
	if err := r.db.FindOne(qCtx, filter).Decode(foundVideo); err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return true, r.logger.CriticalPropagate(err)
	}

	if !foundVideo.ID.Value.IsZero() && foundVideo.ID.Value != video.ID.Value {
		return true, nil
	}
	return false, nil
}

func (r *VideoRepository) Remove(ctx context.Context, video *agg.Video) error {
	qCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	res, err := r.db.DeleteOne(qCtx, bson.M{"_id": video.ID.Value})
	if err != nil {
		return r.logger.ErrorPropagate(err)
	}

	if res.DeletedCount == 0 { // checking the video is really deleted
		return r.logger.CriticalPropagate(errors.New("video with id " + video.ID.Value.Hex() + " was not deleted"))
	}

	return nil
}
