package server

type ServerConfig struct {
	ServerPort string `envconfig:"PORT" default:"3300"`
}
