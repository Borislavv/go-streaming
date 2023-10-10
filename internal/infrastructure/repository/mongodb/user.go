package mongodb

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
	"time"
)

const UserCollection = "users"

var (
	UserNotFoundByIdError    = errors.NewEntityNotFoundError("user", "id")
	UserNotFoundByEmailError = errors.NewEntityNotFoundError("user", "email")
	UserInsertingFailedError = errors.NewInternalValidationError("unable to store 'user' or get inserted 'id'")
	UserWasNotDeletedError   = errors.NewInternalValidationError("user was not deleted")
)

type UserRepository struct {
	db      *mongo.Collection
	mu      *sync.Mutex
	logger  logger.Logger
	timeout time.Duration
}

func NewUserRepository(db *mongo.Database, logger logger.Logger, timeout time.Duration) *UserRepository {
	return &UserRepository{
		db:      db.Collection(UserCollection),
		logger:  logger,
		mu:      &sync.Mutex{},
		timeout: timeout,
	}
}

func (r *UserRepository) Find(ctx context.Context, id vo.ID) (user *agg.User, err error) {
	qCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	user = &agg.User{}
	if err = r.db.FindOne(qCtx, bson.M{"_id": bson.M{"$eq": id.Value}}).Decode(user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, r.logger.InfoPropagate(UserNotFoundByIdError)
		}
		return nil, r.logger.ErrorPropagate(err)
	}

	return user, nil
}
