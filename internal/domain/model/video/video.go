package video

type Video struct {
	path string
}

func New(path string) *Video {
	return &Video{path: path}
}

func (v *Video) GetPath() string {
	return v.path
}
