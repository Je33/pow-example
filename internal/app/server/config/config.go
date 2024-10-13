package config

type Config struct {
	FilePath      string `envconfig:"FILE_PATH" default:"quotes.json"`
	ServerAddress string `envconfig:"SERVER_ADDRESS" default:":8080"`
	ServerNetwork string `envconfig:"SERVER_NETWORK" default:"tcp"`
	MaxClients    int64  `envconfig:"MAX_CLIENTS" default:"10"`
	Difficulty    int64  `envconfig:"DIFFICULTY" default:"5"`
	LogLevel      string `envconfig:"LOG_LEVEL" default:"debug"`
}
