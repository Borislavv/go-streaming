package mongodb

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (r *VideoRepository) Insert(ctx context.Context, video *agg.Video) (string, error) {
	qCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	res, err := r.db.InsertOne(qCtx, video, options.InsertOne())
	if err != nil {
		return "", err
	}

	if id, ok := res.InsertedID.(primitive.ObjectID); ok {
		return id.Hex(), nil
	}

	return "", nil
}

func (r *VideoRepository) InsertMany(ctx context.Context, videos []agg.Video) error {
	return nil
}
