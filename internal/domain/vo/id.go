package vo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ID struct {
	Value primitive.ObjectID `json:"value" bson:"_id,omitempty"`
}
