package stream

type config struct {
	// server
	Host      string `env:"RESOURCES_SERVER_HOST" envDefault:"0.0.0.0"`
	Port      string `env:"RESOURCES_SERVER_PORT" envDefault:"9988"`
	Transport string `env:"RESOURCES_SERVER_TRANSPORT_PROTOCOL" envDefault:"tcp"`
	// database
	MongoUri string `env:"MONGO_URI" envDefault:"mongodb://database:27017/streaming"`
}
