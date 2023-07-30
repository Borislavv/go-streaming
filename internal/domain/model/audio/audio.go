package audio

type Audio struct {
	path string
}

func New(path string) *Audio {
	return &Audio{path: path}
}

func (a *Audio) GetPath() string {
	return a.path
}
