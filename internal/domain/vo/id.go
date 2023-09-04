package vo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ID struct {
	Value primitive.ObjectID `json:"id" bson:"_id,omitempty"`
}
