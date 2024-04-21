package uploader

import (
	dtointerface "github.com/Borislavv/video-streaming/internal/domain/dto/interface"
	"github.com/Borislavv/video-streaming/internal/domain/errtype"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	fileinterface "github.com/Borislavv/video-streaming/internal/infrastructure/service/uploader/file/interface"
)

const MultipartFormUploadingType = "multipart_form"

// MultipartFormUploader is a service which represents functionality
// for uploader a full file from *http.Request into fileStorage.
// This approach of uploading takes a much more RAM but works more fast than MultipartPartUploader.
// If you care of performance, you need use this approach, but take care of using RAM and set up
// the appropriate value of 'inMemoryFileSizeThreshold' through env. configuration.
type MultipartFormUploader struct {
	logger                    loggerinterface.Logger
	fileStorage               fileinterface.Storage
	fileNameComputer          fileinterface.NameComputer
	formFilename              string
	maxFilesize               int64
	inMemoryFileSizeThreshold int64
}

func NewNativeUploader(serviceContainer diinterface.ServiceContainer) (*MultipartFormUploader, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	storageService, err := serviceContainer.GetFileStorageService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	fileNameComputer, err := serviceContainer.GetFileNameComputerService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	config, err := serviceContainer.GetConfig()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return &MultipartFormUploader{
		logger:                    loggerService,
		fileStorage:               storageService,
		fileNameComputer:          fileNameComputer,
		formFilename:              config.ResourceFormFilename,
		inMemoryFileSizeThreshold: config.ResourceInMemoryFileSizeThreshold,
	}, nil
}

// Upload method will be store a file on the disk and calculate a new hashed name. Request DTO mutation!
func (u *MultipartFormUploader) Upload(reqDTO dtointerface.UploadResourceRequest) (err error) {
	// request will be parsed and stored in the memory if it is under the RAM threshold,
	// otherwise last parts of parsed file will be stored in the tmp files on the disk space
	if err = reqDTO.GetRequest().ParseMultipartForm(u.inMemoryFileSizeThreshold); err != nil {
		return u.logger.LogPropagate(err)
	}

	// receiving a file and header from multipart/form-data
	// by requested fileNameComputer which is stored in the `formFilename` const.
	formFile, header, err := reqDTO.GetRequest().FormFile(u.formFilename)
	if err != nil {
		return u.logger.LogPropagate(err)
	}
	defer func() { _ = formFile.Close() }()

	// TODO must be added filesize for check uniqueness
	computedFilename, err := u.fileNameComputer.Get(
		reqDTO.GetUserID(),
		header.Filename,
		header.Header.Get("Content-Type"),
		header.Header.Get("Content-Disposition"),
	)

	// checking whether the being uploaded resource already exists
	has, err := u.fileStorage.Has(reqDTO.GetUserID(), computedFilename)
	if err != nil {
		return u.logger.LogPropagate(err)
	}
	if has { // if being uploading resource is already exists, then throw an error
		return u.logger.LogPropagate(errtype.NewResourceAlreadyExistsError(header.Filename))
	}

	// saving a file on disk and calculating new hashed name with full qualified path
	length, filename, filepath, err := u.fileStorage.Store(reqDTO.GetUserID(), computedFilename, formFile)
	if err != nil {
		return u.logger.LogPropagate(err)
	}

	// mutate request reqDTO
	reqDTO.SetOriginFilename(header.Filename)
	reqDTO.SetUploadedFilename(filename)
	reqDTO.SetUploadedFilepath(filepath)
	reqDTO.SetUploadedFilesize(length)
	reqDTO.SetUploadedFiletype(header.Header.Get("Content-Type"))

	return nil
}
