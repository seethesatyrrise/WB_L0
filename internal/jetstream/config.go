package jetstream

type JSConfig struct {
	URL string `default:"localhost:4222" envconfig:"URL"`
}
