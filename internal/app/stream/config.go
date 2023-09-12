package stream

type config struct {
	// server
	Host      string `env:"STREAMING_SERVER_HOST" envDefault:"0.0.0.0"`
	Port      string `env:"STREAMING_SERVER_PORT" envDefault:"9988"`
	Transport string `env:"STREAMING_SERVER_TRANSPORT_PROTOCOL" envDefault:"tcp"`
	// database
	MongoUri string `env:"MONGO_URI" envDefault:"mongodb://mongodb:27017/streaming"`
	MongoDb  string `env:"MONGO_DATABASE" envDefault:"streaming"`
}
