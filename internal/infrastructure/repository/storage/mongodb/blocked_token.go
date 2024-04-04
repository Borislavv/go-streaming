package mongodb

import (
	"context"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/errtype"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

const BlockedTokensCollection = "blockedTokens"

type BlockedTokenRepository struct {
	db      *mongo.Collection
	mu      *sync.Mutex
	logger  loggerinterface.Logger
	timeout time.Duration
}

func NewBlockedTokenRepository(serviceContainer diinterface.ContainerManager) (*BlockedTokenRepository, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	mongodb, err := serviceContainer.GetMongoDatabase()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	cfg, err := serviceContainer.GetConfig()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	timeout, err := time.ParseDuration(cfg.MongoTimeout)
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return &BlockedTokenRepository{
		db:      mongodb.Collection(BlockedTokensCollection),
		logger:  loggerService,
		mu:      &sync.Mutex{},
		timeout: timeout,
	}, nil
}

func (r *BlockedTokenRepository) Insert(ctx context.Context, token *agg.BlockedToken) error {
	qCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	res, err := r.db.InsertOne(qCtx, token, options.InsertOne())
	if err != nil {
		return r.logger.ErrorPropagate(err)
	}

	if _, ok := res.InsertedID.(primitive.ObjectID); !ok {
		return r.logger.LogPropagate(
			errtype.NewInternalValidationError(
				fmt.Sprintf("error occurred while inserting a blocked token '%v'", token),
			),
		)
	}

	return nil
}

func (r *BlockedTokenRepository) Has(ctx context.Context, token string) (found bool, err error) {
	qCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	filter := bson.M{"value": token}

	if err = r.db.FindOne(qCtx, filter).Decode(&agg.BlockedToken{}); err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, r.logger.ErrorPropagate(err)
	}

	return true, nil
}
