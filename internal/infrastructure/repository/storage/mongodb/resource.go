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

const ResourcesCollection = "resources"

var (
	ResourceNotFoundByIdError    = errors.NewEntityNotFoundError("resource", "id")
	ResourceInsertingFailedError = errors.NewInternalRepositoryError("unable to store 'resource' or retrieve inserted 'id'")
)

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

func (r *ResourceRepository) FindOneByID(ctx context.Context, q query.FindOneResourceByID) (*agg.Resource, error) {
	qCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	filter := bson.M{
		"_id":    bson.M{"$eq": q.GetID().Value},
		"userID": bson.M{"$eq": q.GetUserID().Value},
	}

	resourceAgg := &agg.Resource{}
	if err := r.db.FindOne(qCtx, filter).Decode(resourceAgg); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, r.logger.InfoPropagate(ResourceNotFoundByIdError)
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
		q := dto.NewResourceGetRequestDTO(vo.ID{Value: oid}, resource.UserID)
		return r.FindOneByID(qCtx, q)
	}

	return nil, r.logger.CriticalPropagate(ResourceInsertingFailedError)
}

func (r *ResourceRepository) Remove(ctx context.Context, resource *agg.Resource) error {
	qCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	res, err := r.db.DeleteOne(qCtx, bson.M{"_id": resource.ID.Value})
	if err != nil {
		return r.logger.ErrorPropagate(err)
	}

	if res.DeletedCount == 0 { // checking the resource was deleted
		return r.logger.CriticalPropagate(UserWasNotDeletedError)
	}

	return nil
}
