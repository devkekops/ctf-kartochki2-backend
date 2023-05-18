package config

type Config struct {
	ServerAddress  string `env:"SERVER_ADDRESS" envDefault:"0.0.0.0:80"`
	DatabaseDSN    string `env:"DATABASE_DSN" envDefault:"../db/sqlite.db"`
	SecretKey      string `env:"SECRET_KEY" envDefault:"hdkH74hEY45ta983!dhf0"`
	LicenseKeyHash string `env:"LICENSE_KEY_HASH" envDefault:"0a821c6e9883ff5ad03ce6e3b615bd12503ad8a48fa82c63059b6ba9106af823"`
}
