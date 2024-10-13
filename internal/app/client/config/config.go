package config

type Config struct {
	ServerAddress string `envconfig:"SERVER_ADDRESS" default:"127.0.0.1:8080"`
	ServerNetwork string `envconfig:"SERVER_NETWORK" default:"tcp"`
	LogLevel      string `envconfig:"LOG_LEVEL" default:"debug"`
	Difficulty    int64  `envconfig:"DIFFICULTY" default:"5"`
}
