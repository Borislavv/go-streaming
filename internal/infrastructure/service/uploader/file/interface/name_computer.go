package file_interface

type NameComputer interface {
	Get(remoteName string, contentType string, contentDisposition string) (filename string, err error)
}
