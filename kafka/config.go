package kafka

type Config struct {
	BootstrapServers []string
}

type ConsumerConfig struct {
	GroupID string
	Topic   string
}
