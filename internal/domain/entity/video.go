package entity

type Video struct {
	Name        string `json:"name" bson:"name"`
	Path        string `json:"path" bson:"path"`
	Description string `json:"description" bson:"description,omitempty"`
}

func NewVideo(name string, path string, description string) *Video {
	return &Video{
		Name:        name,
		Path:        path,
		Description: description,
	}
}

func (v *Video) GetPath() string {
	return v.Path
}
