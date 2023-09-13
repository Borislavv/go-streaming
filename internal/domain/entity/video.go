package entity

type Video struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description,omitempty"`
}

func NewVideo(name string, description string) *Video {
	return &Video{
		Name:        name,
		Description: description,
	}
}
