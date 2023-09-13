package mongodb

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/errs"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

const ResourcesCollection = "resources"

type ResourceRepository struct {
	db      *mongo.Collection
	mu      *sync.Mutex
	logger  logger.Logger
	timeout time.Duration
}

func NewResourceRepository(db *mongo.Database, logger logger.Logger, timeout time.Duration) *ResourceRepository {
	return &ResourceRepository{
		db:      db.Collection(ResourcesCollection),
		logger:  logger,
		mu:      &sync.Mutex{},
		timeout: timeout,
	}
}

func (r *ResourceRepository) Find(ctx context.Context, id vo.ID) (*agg.Resource, error) {
	qCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	resourceAgg := &agg.Resource{}
	if err := r.db.FindOne(qCtx, bson.M{"_id": bson.M{"$eq": id.Value}}).Decode(resourceAgg); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, r.logger.InfoPropagate(errs.NewNotFoundError("resource"))
		}
		return nil, r.logger.ErrorPropagate(err)
	}

	return resourceAgg, nil
}

func (r *ResourceRepository) Insert(ctx context.Context, resource *agg.Resource) (*agg.Resource, error) {
	qCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	res, err := r.db.InsertOne(qCtx, resource, options.InsertOne())
	if err != nil {
		return nil, r.logger.ErrorPropagate(err)
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		return r.Find(qCtx, vo.ID{Value: oid})
	}

	return nil, r.logger.CriticalPropagate("unable to store 'resource' or retrieve inserted 'id'")
}
