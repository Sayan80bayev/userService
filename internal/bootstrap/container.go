package bootstrap

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"userService/internal/config"
	"userService/internal/messaging"
	"userService/internal/model"
	"userService/internal/pkg/storage"
	"userService/internal/repository"
	"userService/pkg/logging"
)

type Container struct {
	DB                 *gorm.DB
	Redis              *redis.Client
	Minio              *minio.Client
	FileService        storage.FileService
	Producer           messaging.Producer
	Consumer           messaging.Consumer
	Config             *config.Config
	UserRepositoryImpl *repository.UserRepositoryImpl
}

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

	minioClient := storage.Init(cfg)
	fileService := storage.NewMinioStorage(minioClient, cfg)

	producer, err := messaging.NewKafkaProducer("localhost:9092", "posts-events")
	if err != nil {
		logger.Fatal("Error creating Kafka Producer:", err)
		return nil, err
	}

	userRepo := repository.NewUserRepository(db)

	consumer, err := initKafkaConsumer(redisClient, userRepo, fileService)
	if err != nil {
		return nil, err
	}

	logger.Info("✅ Dependencies initialized successfully")

	return &Container{
		DB:                 db,
		Redis:              redisClient,
		Minio:              minioClient,
		FileService:        fileService,
		Producer:           producer,
		Consumer:           consumer,
		Config:             cfg,
		UserRepositoryImpl: userRepo,
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

func initKafkaConsumer(redisClient *redis.Client, userRepo messaging.CacheUserRepository, fileService storage.FileService) (messaging.Consumer, error) {
	consumer, err := messaging.NewConsumer(
		messaging.ConsumerConfig{
			BootstrapServers: "localhost:9092",
			GroupID:          "post-group",
			Topics:           []string{"posts-events"},
		},
		redisClient,
		userRepo,
		fileService,
	)

	if err != nil {
		logging.GetLogger().Fatal("Error initializing Kafka Consumer:", err)
		return nil, err
	}

	return consumer, nil
}
