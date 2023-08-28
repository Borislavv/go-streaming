package entity

type Audio struct {
	path string
}

func NewAudio(path string) *Audio {
	return &Audio{path: path}
}

func (a *Audio) GetPath() string {
	return a.path
}
