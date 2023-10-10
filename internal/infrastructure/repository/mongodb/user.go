package mongodb

import "github.com/Borislavv/video-streaming/internal/domain/errors"

const UserCollection = "users"

var (
	UserNotFoundByIdError    = errors.NewEntityNotFoundError("user", "id")
	UserNotFoundByEmailError = errors.NewEntityNotFoundError("user", "email")
	UserInsertingFailedError = errors.NewInternalValidationError("unable to store 'user' or get inserted 'id'")
	UserWasNotDeletedError   = errors.NewInternalValidationError("user was not deleted")
)
