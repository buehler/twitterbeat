package beater

type TwitterConfig struct {
	Period  *int64
	Twitter struct {
		Names          *[]string
		AccessKey      *string `yaml:"access_key"`
		AccessSecret   *string `yaml:"access_secret"`
		ConsumerKey    *string `yaml:"consumer_key"`
		ConsumerSecret *string `yaml:"consumer_secret"`
	}
}

type TwitterConfigYaml struct {
	Input TwitterConfig
}
