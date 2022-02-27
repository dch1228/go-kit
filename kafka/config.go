package kafka

type Config struct {
	BootstrapServers []string `mapstructure:"bootstrap-servers"`
}

type ConsumerConfig struct {
	GroupID string `mapstructure:"group-id"`
	Topic   string `mapstructure:"topic"`
}
