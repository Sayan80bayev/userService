package config

import (
	"github.com/Sayan80bayev/go-project/pkg/logging"
	"github.com/spf13/viper"
)

type Config struct {
	MongoURI    string `mapstructure:"MONGO_URI"`
	MongoDBName string `mapstructure:"MONGO_DB_NAME"`
	Port        string `mapstructure:"PORT"`
	GrpcPort    string `mapstructure:"GRPC_PORT"`

	AccessKey   string `mapstructure:"ACCESS_KEY"`
	SecretKey   string `mapstructure:"SECRET_KEY"`
	MinioBucket string `mapstructure:"MINIO_BUCKET"`
	MinioHost   string `mapstructure:"MINIO_HOST"`
	MinioPort   string `mapstructure:"MINIO_PORT"`

	RedisAddr string `mapstructure:"REDIS_ADDR"`
	RedisPass string `mapstructure:"REDIS_PASS"`

	KafkaBrokers        []string `mapstructure:"KAFKA_BROKERS"`
	KafkaTopic          string   `mapstructure:"KAFKA_TOPIC"`
	KafkaProducerTopic  string   `mapstructure:"KAFKA_PRODUCER_TOPIC"`
	KafkaConsumerGroup  string   `mapstructure:"KAFKA_CONSUMER_GROUP"`
	KafkaConsumerTopics []string `mapstructure:"KAFKA_CONSUMER_TOPICS"`
	KeycloakURL         string   `mapstructure:"KEYCLOAK_URL"`
	KeycloakRealm       string   `mapstructure:"KEYCLOAK_REALM"`
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
