package jetstream

type JSConfig struct {
	URL          string `default:"localhost:4222" envconfig:"URL"`
	StreamName   string `default:"MSG"`
	ConsumerName string `default:"server1"`
	SubjectName  string `default:"msg"`
}
