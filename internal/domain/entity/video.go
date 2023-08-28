package entity

type Video struct {
	Name        string `bson:"name"`
	Path        string `bson:"path"`
	Description string `bson:"description,omitempty"`
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
