package agg

import (
	"github.com/Borislavv/video-streaming/data/vo"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Video struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Video     entity.Video       `bson:",inline"`
	Timestamp vo.Timestamp       `bson:",inline"`
}
