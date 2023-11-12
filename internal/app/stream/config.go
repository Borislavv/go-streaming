package stream

type config struct {
	// >>> SERVER <<<
	// Host is an WebSocket server serving host.
	Host string `env:"STREAMING_SERVER_HOST" envDefault:"0.0.0.0"`
	// Port is an WebSocket server serving port.
	Port string `env:"STREAMING_SERVER_PORT" envDefault:"9988"`
	// Transport is an WebSocket server transport protocol.
	// If you are not concerned about the loss part of packets and this is not a problem for you, then use the UDP,
	// because this will give you a performance gain (due to the server will not check of packages number and them ordering).
	// Otherwise, if your data needs to be in safe, and you cannot afford to lose it, use the TCP.
	Transport string `env:"STREAMING_SERVER_TRANSPORT_PROTOCOL" envDefault:"tcp" opts:"tcp,udp"`
	// >>> DATABASE <<<
	// MongoUri is a simple MongoDb DSN string for connect to database.
	MongoUri string `env:"MONGO_URI" envDefault:"mongodb://mongodb:27017/streaming"`
	// MongoDb is a name of database into the MongoDb.
	MongoDb string `env:"MONGO_DATABASE" envDefault:"streaming"`
	// >>> FILE READER <<<
	// ChunkSize is a value which means the size of one chunk while reading the file when streaming a resource.
	// By default, it's 1mb.
	ChunkSize int `env:"FILE_READER_CHUNK_SIZE" envDefault:"1048576"`
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
}
