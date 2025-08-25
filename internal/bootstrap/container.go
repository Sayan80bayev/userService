package bootstrap

import (
	"context"
	"fmt"
	"github.com/Sayan80bayev/go-project/pkg/caching"
	"github.com/Sayan80bayev/go-project/pkg/logging"
	"github.com/Sayan80bayev/go-project/pkg/messaging"
	storage "github.com/Sayan80bayev/go-project/pkg/objectStorage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"userService/internal/config"
	"userService/internal/events"
	"userService/internal/repository"
	"userService/internal/service"
)

// Container holds all dependencies
type Container struct {
	DB             *mongo.Database
	Redis          caching.CacheService
	FileStorage    storage.FileStorage
	Producer       messaging.Producer
	Consumer       messaging.Consumer
	UserRepository service.UserRepository
	Config         *config.Config
	JWKSUrl        string
}

// Init initializes all dependencies and returns a container
func Init() (*Container, error) {
	logger := logging.GetLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	db, err := initMongoDatabase(cfg)
	if err != nil {
		return nil, err
	}

	cacheService, err := initRedis(cfg)
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

	userRepository := repository.NewUserRepository(db)

	consumer, err := initKafkaConsumer(cfg, fileStorage, userRepository)
	if err != nil {
		return nil, err
	}

	jwksURL := buildJWKSURL(cfg)

	logger.Info("✅ Dependencies initialized successfully")
	// Wait for shutdown signal

	return &Container{
		DB:             db,
		Redis:          cacheService,
		FileStorage:    fileStorage,
		Producer:       producer,
		Consumer:       consumer,
		Config:         cfg,
		JWKSUrl:        jwksURL,
		UserRepository: userRepository,
	}, nil
}

// --- Helpers ---

func initMongoDatabase(cfg *config.Config) (*mongo.Database, error) {
	logger := logging.GetLogger()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(cfg.MongoURI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		logger.Fatal("Error connecting to MongoDB:", err)
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		logger.Fatal("MongoDB ping failed:", err)
		return nil, err
	}

	logger.Info("Connected to MongoDB")

	return client.Database(cfg.MongoDBName), nil
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

func initKafkaConsumer(cfg *config.Config, fileStorage storage.FileStorage, repo service.UserRepository) (messaging.Consumer, error) {
	consumer, err := messaging.NewKafkaConsumer(messaging.ConsumerConfig{
		BootstrapServers: cfg.KafkaBrokers[0],
		GroupID:          cfg.KafkaConsumerGroup,
		Topics:           cfg.KafkaConsumerTopics,
	})
	if err != nil {
		return nil, fmt.Errorf("kafka consumer init failed: %w", err)
	}

	// Use typed event constants
	consumer.RegisterHandler(events.UserCreated, service.CreateUserHandler(repo))
	consumer.RegisterHandler(events.UserUpdated, service.UserUpdatedHandler(fileStorage))
	consumer.RegisterHandler(events.UserDeleted, service.UserDeletedHandler(fileStorage))

	logging.GetLogger().Infof("Kafka consumer initialized")
	return consumer, nil
}

func buildJWKSURL(cfg *config.Config) string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", cfg.KeycloakURL, cfg.KeycloakRealm)
}
