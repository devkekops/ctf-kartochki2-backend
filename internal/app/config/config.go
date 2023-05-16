package config

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"0.0.0.0:80"`
	SecretKey     string `env:"SECRET_KEY" envDefault:"hdkH74hEY45ta983!dhf0"`
}
