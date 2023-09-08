package mongodb

import (
	"context"
	"errors"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/errs"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"regexp"
	"sync"
	"time"
)

const VideoCollection = "videos"

type VideoRepository struct {
	db      *mongo.Collection
	mu      *sync.Mutex
	buf     []interface{}
	timeout time.Duration
}

func NewVideoRepository(db *mongo.Database, timeout time.Duration) *VideoRepository {
	return &VideoRepository{
		db:      db.Collection(VideoCollection),
		mu:      &sync.Mutex{},
		buf:     []interface{}{},
		timeout: timeout,
	}
}

func (r *VideoRepository) Find(ctx context.Context, id vo.ID) (*agg.Video, error) {
	qCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	videoAgg := &agg.Video{}
	if err := r.db.FindOne(qCtx, bson.M{"_id": bson.M{"$eq": id.Value}}).Decode(videoAgg); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errs.NewNotFoundError("video")
		}
		return nil, err
	}

	return videoAgg, nil
}

func (r *VideoRepository) FindList(ctx context.Context, query dto.ListRequest) ([]*agg.Video, error) {
	qCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	filter := bson.M{}

	if query.GetName() != "" {
		filter["name"] = primitive.Regex{Pattern: query.GetName(), Options: "i"}
	}
	if query.GetPath() != "" {
		filter["path"] = primitive.Regex{Pattern: "^" + regexp.QuoteMeta(query.GetPath()), Options: "i"}
	}

	opts := options.Find().
		SetSkip((int64(query.GetPage()) - 1) * int64(query.GetLimit())).
		SetLimit(int64(query.GetLimit()))

	videos := []*agg.Video{}
	cursor, err := r.db.Find(qCtx, filter, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return videos, nil
		}
		return nil, err
	}

	if err = cursor.All(qCtx, &videos); err != nil {
		return nil, err
	}

	return videos, nil
}

func (r *VideoRepository) Insert(ctx context.Context, video *agg.Video) (*agg.Video, error) {
	qCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	res, err := r.db.InsertOne(qCtx, video, options.InsertOne())
	if err != nil {
		return nil, err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		return r.Find(qCtx, vo.ID{Value: oid})
	}

	return nil, errors.New("unable to store 'video' or retrieve inserted 'id'")
}

func (r *VideoRepository) Update(ctx context.Context, video *agg.Video) (*agg.Video, error) {
	qCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	res, err := r.db.UpdateByID(qCtx, video.ID.Value, bson.M{"$set": video})
	if err != nil {
		return nil, err
	}

	// check the record is really updated
	if res.ModifiedCount > 0 {
		return r.Find(qCtx, video.ID)
	}

	// if changes is not exists, then return the original data
	return video, nil
}

func (r *VideoRepository) Remove(ctx context.Context, video *agg.Video) error {
	qCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	res, err := r.db.DeleteOne(qCtx, bson.M{"_id": video.ID.Value})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 { // checking the video is really deleted
		return errors.New("video with id " + video.ID.Value.Hex() + " was not deleted")
	}

	return nil
}
