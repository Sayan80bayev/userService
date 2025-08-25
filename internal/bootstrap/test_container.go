package bootstrap

import (
	"fmt"
	"github.com/Sayan80bayev/go-project/pkg/messaging"
	"userService/internal/config"
	"userService/internal/repository"
)

func NewTestContainer(mongoURI, kafkaAddr, minioHost, minioPort, redisAddr, jwksURL string) *Container {
	cfg := &config.Config{
		MongoURI:            mongoURI,
		MongoDBName:         "testdb",
		RedisAddr:           redisAddr,
		RedisPass:           "",
		MinioBucket:         "test-bucket",
		MinioHost:           minioHost,
		MinioPort:           minioPort, // usually testcontainer default
		AccessKey:           "admin",
		SecretKey:           "admin123",
		KafkaBrokers:        []string{kafkaAddr},
		KafkaProducerTopic:  "user-events",
		KafkaConsumerGroup:  "user-service-test",
		KafkaConsumerTopics: []string{"user-events"},
	}

	// Mongo
	db, err := initMongoDatabase(cfg)
	if err != nil {
		panic(err)
	}

	// Redis
	cacheService, err := initRedis(cfg)
	if err != nil {
		panic(err)
	}

	// MinIO
	fs, err := initMinio(cfg)
	if err != nil {
		panic(err)
	}

	// Kafka Producer
	producer, err := messaging.NewKafkaProducer(cfg.KafkaBrokers[0], cfg.KafkaProducerTopic)
	if err != nil {
		panic(fmt.Errorf("failed to create Kafka producer: %w", err))
	}

	userRepository := repository.NewUserRepository(db)

	// Kafka Consumer
	consumer, err := initKafkaConsumer(cfg, fs, userRepository)
	if err != nil {
		panic(err)
	}
	// Use typed event constants

	return &Container{
		DB:             db,
		Redis:          cacheService,
		FileStorage:    fs,
		Producer:       producer,
		Consumer:       consumer,
		UserRepository: userRepository,
		Config:         cfg,
		JWKSUrl:        jwksURL,
	}
}
