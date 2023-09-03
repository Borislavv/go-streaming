package entity

type Audio struct {
	Name        string `bson:"name"`
	Path        string `bson:"path"`
	Description string `bson:"description,omitempty"`
}

func NewAudio(name string, path string, description string) *Audio {
	return &Audio{
		Name:        name,
		Path:        path,
		Description: description,
	}
}

func (v *Audio) GetPath() string {
	return v.Path
}
