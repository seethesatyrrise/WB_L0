package jetstream

type JSConfig struct {
	URL          string `default:":4222" envconfig:"URL"`
	StreamName   string `default:"MSG"`
	ConsumerName string `default:"server11"`
	SubjectName  string `default:"msg"`
}
