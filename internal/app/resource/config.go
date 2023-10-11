package resource

type config struct {
	// >>> API <<<
	// ApiVersionPrefix is a value which will be used as your RestAPI controllers version prefix.
	// For example: {{schema}}://{{host}}:{{port}}{{ApiVersionPrefix}}/{{additionalControllerPath}}
	ApiVersionPrefix string `env:"API_VERSION_PREFIX" envDefault:"/api/v1"`
	// RenderVersionPrefix is a value which will be used as your native rendering controllers version prefix.
	// For example: {{schema}}://{{host}}:{{port}}{{RenderVersionPrefix}}/{{additionalControllerPath}}
	// By default it's an empty string.
	RenderVersionPrefix string `env:"RENDER_VERSION_PREFIX" envDefault:""`
	// StaticVersionPrefix is a value which will be used as your static files controllers version prefix.
	// For example: {{schema}}://{{host}}:{{port}}{{StaticVersionPrefix}}/{{additionalControllerPath}}
	// By default it's an empty string.
	StaticVersionPrefix string `env:"STATIC_VERSION_PREFIX" envDefault:""`
	// >>> SERVER <<<
	// Host is an HTTP server serving host.
	Host string `env:"RESOURCES_SERVER_HOST" envDefault:"0.0.0.0"`
	// Port is an HTTP server serving port.
	Port string `env:"RESOURCES_SERVER_PORT" envDefault:"8000"`
	// Transport is an HTTP server transport protocol.
	// If you are not concerned about the loss part of packets and this is not a problem for you, then use the UDP,
	// because this will give you a performance gain (due to the server will not check of packages number and them ordering).
	// Otherwise, if your data needs to be in safe, and you cannot afford to lose it, use the TCP.
	Transport string `env:"RESOURCES_SERVER_TRANSPORT_PROTOCOL" envDefault:"tcp" opts:"tcp,udp"`
	// >>> DATABASE <<<
	// MongoUri is a simple MongoDb DSN string for connect to database.
	MongoUri string `env:"MONGO_URI" envDefault:"mongodb://mongodb:27017/streaming"`
	// MongoDb is a name of database into the MongoDb.
	MongoDb string `env:"MONGO_DATABASE" envDefault:"streaming"`
	// >>> APPLICATION <<<
	// UploadingStrategy is an uploading strategy which will be used for upload files on the server.
	// 	1. 'muiltipart_form' is a strategy which used builtin sugar approach. It will be parsing a whole file into the
	//		memory (if a file more than InMemoryFileSizeThreshold, it will be saved on the disk, otherwise, it will be
	//		loaded in the RAM).
	//	2. 'muiltipart_part' is a strategy which used lower level implementation which based on the reading by parts
	//		from raw form data.
	// 	If you care of application performance (speed of uploading directly) and you have enough RAM, then use
	//	the 'muiltipart_form' approach and increase the value of InMemoryFileSizeThreshold variable.
	//	Otherwise, use 'muiltipart_part' because it takes a much lower RAM per file uploading.
	//	For example: for upload the file which weight is 50mb. it will take around 10mb. of your RAM.
	UploadingStrategy string `env:"UPLOADER_TYPE" envDefault:"muiltipart_part" opts:"muiltipart_form,muiltipart_part"`
	// ResourceFormFilename is a value which will be used for extract a file from the form by given string.
	// *Used only with the 'muiltipart_form' strategy because the 'muiltipart_part' will search the first form file.
	//	Be careful and don't send more than one file per request in one form.
	ResourceFormFilename string `env:"RESOURCE_FORM_FILENAME" envDefault:"resource"`
	// MaxFilesizeThreshold is a threshold value which means the max. weight of uploading file in bytes.
	// By default, it's 10gb per file.
	MaxFilesizeThreshold int64 `env:"MAX_UPLOADING_FILESIZE" envDefault:"10000000000"`
	// InMemoryFileSizeThreshold is a threshold value which means the max. weight of uploading file in bytes
	// which may be loaded in the RAM. If file weight is more this value, than it will be loaded on the disk (slow operation).
	// By default, it's 100mb per file.
	InMemoryFileSizeThreshold int64 `env:"IN_MEMORY_FILE_SIZE_THRESHOLD" envDefault:"104857600"`
	// AdminContactEmail is a target administrator contact email address for takes a users errors reports.
	AdminContactEmail string `env:"ADMIN_CONTACT_EMAIL_ADDRESS" envDefault:"glazunov2142@gmail.com"`
	// >>> LOGGER <<<
	// LoggerErrorsBufferCap is errors channel capacity.
	// Logger is basing on the go channels, this value will be sat up as capacity.
	LoggerErrorsBufferCap int `env:"LOGGER_ERRORS_BUFFER_CAPACITY" envDefault:"10"`
	// LoggerRequestsBufferCap is requests channel capacity.
	// Use only when you are logging input requests/responses.
	LoggerRequestsBufferCap int `env:"LOGGER_REQUESTS_BUFFER_CAPACITY" envDefault:"10"`
}
