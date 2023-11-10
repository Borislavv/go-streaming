package mongodb

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"github.com/Borislavv/video-streaming/internal/infrastructure/repository/query"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

const VideosCollection = "videos"

var (
	VideoNotFoundByIdError         = errors.NewEntityNotFoundError("video", "id")
	VideoNotFoundByNameError       = errors.NewEntityNotFoundError("video", "name")
	VideoNotFoundByResourceIdError = errors.NewEntityNotFoundError("video", "resource.id")
	VideoInsertingFailedError      = errors.NewInternalValidationError("unable to store 'video' or get inserted 'id'")
	VideoWasNotDeletedError        = errors.NewInternalValidationError("video was not deleted")
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

func (r *VideoRepository) FindOneByID(ctx context.Context, q query.FindOneVideoByID) (*agg.Video, error) {
	qCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	filter := bson.M{
		"_id":    bson.M{"$eq": q.GetID().Value},
		"userID": bson.M{"$eq": q.GetUserID().Value},
	}

	video := &agg.Video{}
	if err := r.db.FindOne(qCtx, filter).Decode(video); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, r.logger.InfoPropagate(VideoNotFoundByIdError)
		}
		return nil, r.logger.ErrorPropagate(err)
	}

	return video, nil
}

func (r *VideoRepository) FindList(ctx context.Context, q query.FindVideoList) (list []*agg.Video, total int64, err error) {
	qCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	filter := bson.M{"userID": q.GetUserID().Value}

	if q.GetName() != "" {
		filter["name"] = primitive.Regex{Pattern: q.GetName(), Options: "i"}
	}
	if !q.GetCreatedAt().IsZero() {
		y := q.GetCreatedAt().Year()
		m := q.GetCreatedAt().Month()
		d := q.GetCreatedAt().Day()

		filter["createdAt"] = bson.M{
			"$gt":  time.Date(y, m, d, 0, 0, 0, 0, time.UTC),
			"$lte": time.Date(y, m, d, 23, 59, 59, 0, time.UTC),
		}
	} else if !q.GetFrom().IsZero() || !q.GetTo().IsZero() {
		createdAtFilter := bson.M{}
		if !q.GetFrom().IsZero() {
			createdAtFilter["$gt"] = q.GetFrom()
		}
		if !q.GetTo().IsZero() {
			createdAtFilter["$lte"] = q.GetTo()
		}
		filter["createdAt"] = createdAtFilter
	}

	opts := options.Find().
		SetSkip((int64(q.GetPage()) - 1) * int64(q.GetLimit())).
		SetLimit(int64(q.GetLimit()))

	wg := sync.WaitGroup{}
	wg.Add(2)

	list = []*agg.Video{}
	go func() {
		defer wg.Done()

		c, e := r.db.Find(qCtx, filter, opts)
		if e != nil && e != mongo.ErrNoDocuments {
			r.logger.Error(e)
			return
		}
		defer func() { _ = c.Close(qCtx) }()

		if e = c.All(qCtx, &list); e != nil {
			r.logger.Error(e)
		}
	}()

	total = 0
	go func() {
		defer wg.Done()

		c, e := r.db.CountDocuments(qCtx, filter)
		if err != nil {
			r.logger.Error(e)
			return
		}

		total = c
	}()

	wg.Wait()

	return list, total, nil
}

func (r *VideoRepository) FindOneByName(ctx context.Context, q query.FindOneVideoByName) (*agg.Video, error) {
	qCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	filter := bson.M{"name": q.GetName(), "userID": q.GetUserID().Value}

	video := &agg.Video{}
	if err := r.db.FindOne(qCtx, filter).Decode(video); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, VideoNotFoundByNameError
		}
		return nil, r.logger.LogPropagate(err)
	}

	return video, nil
}

func (r *VideoRepository) FindOneByResourceID(ctx context.Context, q query.FindOneVideoByResourceID) (*agg.Video, error) {
	qCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	filter := bson.M{"resource._id": q.GetResourceID().Value, "userID": q.GetUserID().Value}

	video := &agg.Video{}
	if err := r.db.FindOne(qCtx, filter).Decode(video); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, VideoNotFoundByResourceIdError
		}
		return nil, r.logger.LogPropagate(err)
	}

	return video, nil
}

func (r *VideoRepository) Insert(ctx context.Context, video *agg.Video) (*agg.Video, error) {
	qCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	res, err := r.db.InsertOne(qCtx, video, options.InsertOne())
	if err != nil {
		return nil, r.logger.ErrorPropagate(err)
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		return r.FindOneByID(qCtx, dto.NewVideoGetRequestDTO(vo.ID{Value: oid}, video.UserID))
	}

	return nil, r.logger.CriticalPropagate(VideoInsertingFailedError)
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
		return r.FindOneByID(qCtx, dto.NewVideoGetRequestDTO(video.ID, video.UserID))
	}

	// if changes is not exists, then return the original data
	return video, nil
}

func (r *VideoRepository) Remove(ctx context.Context, video *agg.Video) error {
	qCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	res, err := r.db.DeleteOne(qCtx, bson.M{"_id": video.ID.Value})
	if err != nil {
		return r.logger.ErrorPropagate(err)
	}

	if res.DeletedCount == 0 { // checking the video is really deleted
		return r.logger.CriticalPropagate(VideoWasNotDeletedError)
	}

	return nil
}
