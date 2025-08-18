package bootstrap

import (
	"context"
	"fmt"
	"github.com/Sayan80bayev/go-project/pkg/logging"
	"github.com/Sayan80bayev/go-project/pkg/messaging"
	storage "github.com/Sayan80bayev/go-project/pkg/objectStorage"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"userService/internal/config"
	"userService/internal/model"
	"userService/internal/service"
)

// Container is a structure that contains all components for DI.
type Container struct {
	DB          *gorm.DB
	Redis       *redis.Client
	Minio       *minio.Client
	FileStorage storage.FileStorage
	Producer    messaging.Producer
	Consumer    messaging.Consumer
	Config      *config.Config
	JWKSUrl     string
}

// Init initializes container with components.
func Init() (*Container, error) {
	logger := logging.GetLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Error loading configuration:", err)
		return nil, err
	}

	db, err := initDatabase(cfg)
	if err != nil {
		return nil, err
	}

	redisClient, err := initRedis(cfg)
	if err != nil {
		return nil, err
	}

	minioCfg := &storage.MinioConfig{
		Bucket:    cfg.MinioBucket,
		Host:      cfg.MinioHost,
		AccessKey: cfg.AccessKey,
		SecretKey: cfg.SecretKey,
		Port:      cfg.MinioPort,
	}

	fileStorage := storage.GetMinioStorage(minioCfg)

	producer, err := messaging.GetProducer(cfg.KafkaBrokers[0], cfg.KafkaProducerTopic)
	if err != nil {
		logger.Fatal("Error creating Kafka Producer:", err)
		return nil, err
	}

	consumer, err := initKafkaConsumer(cfg, fileStorage)
	if err != nil {
		return nil, err
	}

	logger.Info("✅ Dependencies initialized successfully")

	jwksURL := buildJWKSURL(cfg)

	return &Container{
		DB:          db,
		Redis:       redisClient,
		FileStorage: fileStorage,
		Producer:    producer,
		Consumer:    consumer,
		Config:      cfg,
		JWKSUrl:     jwksURL,
	}, nil
}

func initDatabase(cfg *config.Config) (*gorm.DB, error) {
	logger := logging.GetLogger()

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		logger.Fatal("Error connecting to the database:", err)
		return nil, err
	}

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		logger.Fatal("Database migration error:", err)
		return nil, err
	}

	return db, nil
}

func initRedis(cfg *config.Config) (*redis.Client, error) {
	logger := logging.GetLogger()

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
		DB:       0,
	})

	ctx := context.Background()
	if _, err := client.Ping(ctx).Result(); err != nil {
		logger.Fatal("Error connecting to Redis:", err)
		return nil, err
	}

	return client, nil
}

func initKafkaConsumer(cfg *config.Config, fileStorage storage.FileStorage) (messaging.Consumer, error) {
	consumer, err := messaging.GetConsumer(
		messaging.ConsumerConfig{
			BootstrapServers: cfg.KafkaBrokers[0],
			GroupID:          cfg.KafkaConsumerGroup,
			Topics:           cfg.KafkaConsumerTopics,
		},
	)

	if err != nil {
		logging.GetLogger().Fatal("Error initializing Kafka Consumer:", err)
		return nil, err
	}

	consumer.RegisterHandler("UserUpdated", service.UserUpdatedHandler(fileStorage))
	consumer.RegisterHandler("UserDeleted", service.UserDeletedHandler(fileStorage))

	return consumer, nil
}

func buildJWKSURL(cfg *config.Config) string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", cfg.KeycloakURL, cfg.KeycloakRealm)
}
