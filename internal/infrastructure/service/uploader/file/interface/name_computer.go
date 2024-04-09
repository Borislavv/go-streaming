package fileinterface

import "github.com/Borislavv/video-streaming/internal/domain/vo"

type NameComputer interface {
	Get(userID vo.ID, remoteName string, contentType string, contentDisposition string) (filename string, err error)
}
