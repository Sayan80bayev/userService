package bootstrap

import (
	"fmt"
	"github.com/Sayan80bayev/go-project/pkg/caching"
	"github.com/Sayan80bayev/go-project/pkg/logging"
	"github.com/Sayan80bayev/go-project/pkg/messaging"
	storage "github.com/Sayan80bayev/go-project/pkg/objectStorage"
	"github.com/minio/minio-go/v7"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"userService/internal/config"
	"userService/internal/events"
	"userService/internal/model"
	"userService/internal/service"
)

// Container holds all dependencies
type Container struct {
	DB          *gorm.DB
	Redis       caching.CacheService
	Minio       *minio.Client
	FileStorage storage.FileStorage
	Producer    messaging.Producer
	Consumer    messaging.Consumer
	Config      *config.Config
	JWKSUrl     string
}

// Init initializes all dependencies and returns a container
func Init() (*Container, error) {
	logger := logging.GetLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	db, err := initDatabase(cfg)
	if err != nil {
		return nil, err
	}

	redisClient, err := initRedis(cfg)
	if err != nil {
		return nil, err
	}

	fileStorage, err := initMinio(cfg)
	if err != nil {
		return nil, err
	}

	producer, err := messaging.NewKafkaProducer(cfg.KafkaBrokers[0], cfg.KafkaProducerTopic)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	consumer, err := initKafkaConsumer(cfg, fileStorage)
	if err != nil {
		return nil, err
	}

	jwksURL := buildJWKSURL(cfg)

	logger.Info("✅ Dependencies initialized successfully")

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

// --- Helpers ---

func initDatabase(cfg *config.Config) (*gorm.DB, error) {
	logger := logging.GetLogger()

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("db connection failed: %w", err)
	}

	if err = db.AutoMigrate(&model.User{}); err != nil {
		return nil, fmt.Errorf("db migration failed: %w", err)
	}

	logger.Info("Postgres connected & migrated")
	return db, nil
}

func initRedis(cfg *config.Config) (*caching.RedisService, error) {
	logger := logging.GetLogger()
	redisCache, err := caching.NewRedisService(caching.RedisConfig{
		DB:       0,
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
	})

	if err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	logger.Info("Redis connected")
	return redisCache, nil
}

func initMinio(cfg *config.Config) (storage.FileStorage, error) {
	logger := logging.GetLogger()

	minioCfg := &storage.MinioConfig{
		Bucket:    cfg.MinioBucket,
		Host:      cfg.MinioHost,
		AccessKey: cfg.AccessKey,
		SecretKey: cfg.SecretKey,
		Port:      cfg.MinioPort,
	}

	fs, err := storage.NewMinioStorage(minioCfg)
	if err != nil {
		return nil, fmt.Errorf("minio init failed: %w", err)
	}

	logger.Infof("Minio connected: bucket=%s host=%s", cfg.MinioBucket, cfg.MinioHost)
	return fs, nil
}

func initKafkaConsumer(cfg *config.Config, fileStorage storage.FileStorage) (messaging.Consumer, error) {
	consumer, err := messaging.NewKafkaConsumer(messaging.ConsumerConfig{
		BootstrapServers: cfg.KafkaBrokers[0],
		GroupID:          cfg.KafkaConsumerGroup,
		Topics:           cfg.KafkaConsumerTopics,
	})
	if err != nil {
		return nil, fmt.Errorf("kafka consumer init failed: %w", err)
	}

	// Use typed event constants
	consumer.RegisterHandler(events.UserUpdated, service.UserUpdatedHandler(fileStorage))
	consumer.RegisterHandler(events.UserDeleted, service.UserDeletedHandler(fileStorage))

	logging.GetLogger().Info("Kafka consumer initialized")
	return consumer, nil
}

func buildJWKSURL(cfg *config.Config) string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", cfg.KeycloakURL, cfg.KeycloakRealm)
}
