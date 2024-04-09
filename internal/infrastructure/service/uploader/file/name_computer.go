package file

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"path/filepath"
)

type NameComputerService struct {
}

func NewNameComputerService() *NameComputerService {
	return &NameComputerService{}
}

// Get - will return computed filename with extension.
func (s *NameComputerService) Get(
	userID vo.ID,
	remoteFilename string,
	contentType string,
	contentDisposition string,
) (filename string, err error) {
	hash := sha256.New()
	if _, err = hash.Write(
		[]byte(
			fmt.Sprintf(
				"%v%v%v%+v",
				userID.Value.Hex(),
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
