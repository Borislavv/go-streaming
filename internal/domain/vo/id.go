package vo

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

type ID struct {
	Value primitive.ObjectID `json:"value" bson:"_id,omitempty"`
}

func (id *ID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%v\"", id.Value.Hex())), nil
}

func (id *ID) UnmarshalJSON(b []byte) error {
	hex := strings.ReplaceAll(string(b), "\"", "")

	oid, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		return err
	}
	id.Value = oid

	return nil
}
