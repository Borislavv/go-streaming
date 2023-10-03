package file

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"
)

type NameService struct {
}

func NewNameService() *NameService {
	return &NameService{}
}

// Get - will return computed filename with extension.
func (s *NameService) Get(
	remoteFilename string,
	contentType string,
	contentDisposition string,
) (filename string, err error) {
	hash := sha256.New()
	if _, err = hash.Write(
		[]byte(
			fmt.Sprintf(
				"%v%v%+v",
				remoteFilename,
				contentType,
				contentDisposition,
			),
		),
	); err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"%v%v",
		hex.EncodeToString(hash.Sum(nil)),
		filepath.Ext(remoteFilename),
	), nil
}
