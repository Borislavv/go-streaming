package uploader

import (
	dtointerface "github.com/Borislavv/video-streaming/internal/domain/dto/interface"
	"github.com/Borislavv/video-streaming/internal/domain/errtype"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	storagerinterface "github.com/Borislavv/video-streaming/internal/domain/service/storager/interface"
	fileinterface "github.com/Borislavv/video-streaming/internal/infrastructure/service/uploader/file/interface"
	"io"
	"mime/multipart"
)

const MultipartPartUploadingType = "multipart_part"

// MultipartPartUploader - is a file ResourceUploadingStrategy which use multipart.Part.
// In such case it takes more time but takes much less memory.
// Approximately, to upload a 50MB file you will need only 10MB of RAM.
type MultipartPartUploader struct {
	logger      loggerinterface.Logger
	storage     storagerinterface.Storage
	filename    fileinterface.NameComputer
	maxFilesize int64
}

func NewPartsUploader(serviceContainer diinterface.ContainerManager) (*MultipartPartUploader, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	storageService, err := serviceContainer.GetFileStorageService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	filenameService, err := serviceContainer.GetFileNameComputerService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return &MultipartPartUploader{
		logger:   loggerService,
		storage:  storageService,
		filename: filenameService,
	}, nil
}

func (u *MultipartPartUploader) Upload(reqDTO dtointerface.UploadResourceRequest) (err error) {
	part, err := u.getFilePart(reqDTO)
	if err != nil {
		return u.logger.LogPropagate(err)
	}

	// TODO must be added filesize for check uniqueness
	computedFilename, err := u.filename.Get(
		reqDTO.GetUserID(),
		part.FileName(),
		part.Header.Get("Content-Type"),
		part.Header.Get("Content-Disposition"),
	)

	// checking whether the being uploaded resource already exists
	has, err := u.storage.Has(reqDTO.GetUserID(), computedFilename)
	if err != nil {
		return u.logger.LogPropagate(err)
	}
	if has { // if being uploading resource is already exists, then throw an error
		return u.logger.LogPropagate(errtype.NewResourceAlreadyExistsError(part.FileName()))
	}

	// saving a file on disk and calculating new hashed name with full qualified path
	length, filename, filepath, err := u.storage.Store(reqDTO.GetUserID(), computedFilename, part)
	if err != nil {
		return u.logger.LogPropagate(err)
	}

	// mutate request reqDTO
	reqDTO.SetOriginFilename(part.FileName())
	reqDTO.SetUploadedFilename(filename)
	reqDTO.SetUploadedFilepath(filepath)
	reqDTO.SetUploadedFilesize(length)
	reqDTO.SetUploadedFiletype(part.Header.Get("Content-Type"))

	return nil
}

func (u *MultipartPartUploader) getFilePart(reqDTO dtointerface.UploadResourceRequest) (part *multipart.Part, err error) {
	// extract the multipart form reader (handling the form as a stream)
	reader, err := reqDTO.GetRequest().MultipartReader()
	if err != nil {
		return nil, u.logger.LogPropagate(err)
	}

	for { // find the part of the form with the target file
		part, err = reader.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, u.logger.LogPropagate(err)
		}

		// check the form part is th target file field
		if part.FileName() != "" {
			return part, nil
		}
	}

	return nil, errtype.NewFormDoesNotContainsUploadedFileError()
}
