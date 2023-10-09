package resource

type config struct {
	// API
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
	// SERVER
	// Host is an HTTP server serving host.
	Host string `env:"RESOURCES_SERVER_HOST" envDefault:"0.0.0.0"`
	// Port is an HTTP server serving port.
	Port string `env:"RESOURCES_SERVER_PORT" envDefault:"8000"`
	// Transport is an HTTP server transport protocol.
	// If you are not concerned about the loss part of packets and this is not a problem for you, then use the UDP,
	// because this will give you a performance gain (due to the server will not check of packages number and them ordering).
	// Otherwise, if your data needs to be in safe, and you cannot afford to lose it, use the TCP.
	Transport string `env:"RESOURCES_SERVER_TRANSPORT_PROTOCOL" envDefault:"tcp" opts:"tcp,udp"`
	// DATABASE
	// MongoUri is a simple MongoDb DSN string for connect to database.
	MongoUri string `env:"MONGO_URI" envDefault:"mongodb://mongodb:27017/streaming"`
	// MongoDb is a name of database into the MongoDb.
	MongoDb string `env:"MONGO_DATABASE" envDefault:"streaming"`
	// application
	Uploader                  string `env:"UPLOADER_TYPE" envDefault:"muiltipart_part"` // supported types: 'muiltipart_form', 'muiltipart_part'
	ResourceFormFilename      string `env:"RESOURCE_FORM_FILENAME" envDefault:"resource"`
	MaxFilesize               int64  `env:"MAX_UPLOADING_FILESIZE" envDefault:"10000000000"`      // 10gb.
	InMemoryFileSizeThreshold int64  `env:"IN_MEMORY_FILE_SIZE_THRESHOLD" envDefault:"104857600"` // 100mb.
	// logger
}
