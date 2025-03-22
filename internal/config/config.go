package config

import (
	"userService/pkg/logging"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL string `mapstructure:"DATABASE_URL"`
	Port        string `mapstructure:"PORT"`
	JWTSecret   string `mapstructure:"JWT_SECRET"`

	AccessKey     string `mapstructure:"ACCESS_KEY"`
	SecretKey     string `mapstructure:"SECRET_KEY"`
	MinioBucket   string `mapstructure:"MINIO_BUCKET"`
	MinioEndpoint string `mapstructure:"MINIO_ENDPOINT"`

	RedisAddr string `mapstructure:"REDIS_ADDR"`
	RedisPass string `mapstructure:"REDIS_PASS"`

	KafkaBrokers        []string `mapstructure:"KAFKA_BROKERS"`
	KafkaTopic          string   `mapstructure:"KAFKA_TOPIC"`
	KafkaProducerTopic  string   `mapstructure:"KAFKA_PRODUCER_TOPIC"`
	KafkaConsumerGroup  string   `mapstructure:"KAFKA_CONSUMER_GROUP"`
	KafkaConsumerTopics []string `mapstructure:"KAFKA_CONSUMER_TOPICS"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile("config/config.yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logging.Instance.Errorf("Couldn't load config.yaml: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
