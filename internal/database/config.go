package database

type DBConfig struct {
	PGAddress  string `envconfig:"ADDRESS" default:":5432"`
	PGUser     string `envconfig:"USER" default:"postgres"`
	PGPassword string `envconfig:"PASSWORD" default:"postgres"`
	PGDatabase string `envconfig:"DATABASE" default:"natsDB"`
}
