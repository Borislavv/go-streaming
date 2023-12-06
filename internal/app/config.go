package resource

type Config struct {
	// >>> RESOURCES HTTP SERVER <<<
	// Host is an HTTP server serving host.
	ResourcesHost string `env:"RESOURCES_SERVER_HOST" envDefault:"0.0.0.0"`
	// Port is an HTTP server serving port.
	ResourcesPort string `env:"RESOURCES_SERVER_PORT" envDefault:"8000"`
	// Transport is an HTTP server transport protocol.
	// If you are not concerned about the loss part of packets and this is not a problem for you, then use the UDP,
	// because this will give you a performance gain (due to the server will not check of packages number and them ord.).
	// Otherwise, if your data needs to be in safe, and you cannot afford to lose it, use the TCP.
	ResourcesTransport string `env:"RESOURCES_SERVER_TRANSPORT_PROTOCOL" envDefault:"tcp" opts:"tcp,udp"`
	// >>> STREAMING WEBSOCKET SERVER <<<
	// Host is an WebSocket server serving host.
	StreamingHost string `env:"STREAMING_SERVER_HOST" envDefault:"0.0.0.0"`
	// Port is an WebSocket server serving port.
	StreamingPort string `env:"STREAMING_SERVER_PORT" envDefault:"9988"`
	// Transport is an WebSocket server transport protocol.
	// If you are not concerned about the loss part of packets and this is not a problem for you, then use the UDP,
	// because this will give you a performance gain (due to the server will not check of packages number and them ordering).
	// Otherwise, if your data needs to be in safe, and you cannot afford to lose it, use the TCP.
	StreamingTransport string `env:"STREAMING_SERVER_TRANSPORT_PROTOCOL" envDefault:"tcp" opts:"tcp,udp"`
	// >>> DATABASE <<<
	// MongoUri is a simple MongoDb DSN string for connect to database.
	MongoUri string `env:"MONGO_URI" envDefault:"mongodb://mongodb:27017/streaming"`
	// MongoDb is a name of database into the MongoDb.
	MongoDb string `env:"MONGO_DATABASE" envDefault:"streaming"`
	// MongoTimeout is a mongo database requests timeout.
	MongoTimeout string `env:"MONGO_TIMEOUT" envDefault:"10s"`
	// >>> APPLICATION <<<
	// JwtSecretSalt is a secret string which further will convert to slice of bytes and will be provided
	// as a salt for signature the jwt tokens.
	// If this variable an empty or omitted then will be generated random salt which is alive while application instance
	// is running (note: you cannot get access to this value, and it will not be provided anywhere as an output. If
	// you need access to this value, and you need to share this value, then sat up your own secret string).
	JwtSecretSalt string `env:"JWT_SECRET_SALT" envDefault:""`
	// JwtTokenIssuer is an issuer of JWT token. This variable helps determine which service issued the token
	// (commonly used for verify that token was created by service of your system, for example,
	// if you have more than one service which able for issue a token).
	JwtTokenIssuer string `env:"JWT_TOKEN_ISSUER" envDefault:"streaming_service"`
	// JwtTokenAcceptedIssuers is a string with another JwtTokenIssuer values separated by delimiter. This values will
	// be accepted while token payload verification.
	JwtTokenAcceptedIssuers string `env:"JWT_TOKEN_ACCEPTED_ISSUERS" envDefault:"auth_service,streaming_service"`
	// JwtTokenExpiresAfter is a value which defined TTL of token in seconds. Default: `86400` (1 day).
	JwtTokenExpiresAfter int64 `env:"JWT_TOKEN_EXPIRES_AFTER" envDefault:"86400"`
	// JwtTokenEncryptAlgo is a value which will be used as encrypt algo for encode the token.
	JwtTokenEncryptAlgo string `env:"JWT_TOKEN_ENCRYPT_ALGO" envDefault:"HS256" opts:"HS256,HS384,HS512"`
	// ResourceUploadingStrategy is an uploading strategy which will be used for upload files on the server.
	// 	1. 'muiltipart_form' is a strategy which used builtin sugar approach. It will be parsing a whole file into the
	//		memory (if a file more than ResourceInMemoryFileSizeThreshold, it will be saved on the disk, otherwise, it will be
	//		loaded in the RAM).
	//	2. 'muiltipart_part' is a strategy which used lower level implementation which based on the reading by parts
	//		from raw form data.
	// 	If you care of application performance (speed of uploading directly) and you have enough RAM, then use
	//	the 'muiltipart_form' approach and increase the value of ResourceInMemoryFileSizeThreshold variable.
	//	Otherwise, use 'muiltipart_part' because it takes a much lower RAM per file uploading.
	//	For example: for upload the file which weight is 50mb. it will take around 10mb. of your RAM.
	ResourceUploadingStrategy string `env:"UPLOADER_TYPE" envDefault:"multipart_part" opts:"multipart_part,multipart_form"`
	// ResourceFormFilename is a value which will be used for extract a file from the form by given string.
	// *Used only with the 'muiltipart_form' strategy because the 'muiltipart_part' will search the first form file.
	//	Be careful and don't send more than one file per request in one form.
	ResourceFormFilename string `env:"RESOURCE_FORM_FILENAME" envDefault:"resource"`
	// ResourceMaxFilesizeThreshold is a threshold value which means the max. weight of uploading file in bytes.
	// By default, it's 5gb per file.
	ResourceMaxFilesizeThreshold int64 `env:"MAX_UPLOADING_FILESIZE" envDefault:"5368709120"`
	// ResourceInMemoryFileSizeThreshold is a threshold value which means the max. weight of uploading file in bytes
	// which may be loaded in the RAM. If file weight is more this value, than it will be loaded on the disk (slow op.).
	// By default, it's 100mb per file.
	ResourceInMemoryFileSizeThreshold int64 `env:"IN_MEMORY_FILE_SIZE_THRESHOLD" envDefault:"104857600"`
	// AdminContactEmail is a target administrator contact email address for takes a users errors reports.
	AdminContactEmail string `env:"ADMIN_CONTACT_EMAIL_ADDRESS" envDefault:"glazunov2142@gmail.com"`
	// >>> API <<<
	// ResourcesApiVersionPrefix is a value which will be used as your RestAPI controllers version prefix.
	// For example: {{schema}}://{{host}}:{{port}}{{ResourcesApiVersionPrefix}}/{{additionalControllerPath}}
	ResourcesApiVersionPrefix string `env:"API_VERSION_PREFIX" envDefault:"/api/v1"`
	// ResourcesRenderVersionPrefix is a value which will be used as your native rendering controllers version prefix.
	// For example: {{schema}}://{{host}}:{{port}}{{ResourcesRenderVersionPrefix}}/{{additionalControllerPath}}
	// By default it's an empty string.
	ResourcesRenderVersionPrefix string `env:"RENDER_VERSION_PREFIX" envDefault:""`
	// ResourcesStaticVersionPrefix is a value which will be used as your static files controllers version prefix.
	// For example: {{schema}}://{{host}}:{{port}}{{ResourcesStaticVersionPrefix}}/{{additionalControllerPath}}
	// By default it's an empty string.
	ResourcesStaticVersionPrefix string `env:"STATIC_VERSION_PREFIX" envDefault:""`
	// >>> LOGGER <<<
	// LoggerErrorsBufferCap is errors channel capacity.
	// Logger is basing on the go channels, this value will be sat up as capacity.
	LoggerErrorsBufferCap int `env:"LOGGER_ERRORS_BUFFER_CAPACITY" envDefault:"10"`
	// LoggerRequestsBufferCap is requests channel capacity.
	// Use only when you are logging input requests/responses.
	LoggerRequestsBufferCap int `env:"LOGGER_REQUESTS_BUFFER_CAPACITY" envDefault:"10"`
	// >>> FILE READER <<<
	// StreamingChunkSize is a value which means the size of one chunk while reading the file when streaming a resource.
	// By default, it's 1mb.
	StreamingChunkSize int `env:"FILE_READER_CHUNK_SIZE" envDefault:"1048576"`
}
