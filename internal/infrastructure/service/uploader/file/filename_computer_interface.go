package file

type NameComputer interface {
	Get(remoteName string, contentType string, contentDisposition string) (filename string, err error)
}
