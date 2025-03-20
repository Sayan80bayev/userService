package messaging

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/redis/go-redis/v9"
	"userService/internal/events"
	"userService/internal/model"
	"userService/internal/pkg/storage"
	"userService/pkg/logging"
)

type CacheUserRepository interface {
	GetUserByUsername(username string) (*model.User, error)
	GetAllUsers() ([]*model.User, error)
}

type ConsumerConfig struct {
	BootstrapServers string
	GroupID          string
	Topics           []string
}

type KafkaConsumer struct {
	config      ConsumerConfig
	consumer    *kafka.Consumer
	redisClient *redis.Client
	userRepo    CacheUserRepository
	fileService storage.FileService
}

func NewConsumer(config ConsumerConfig, redisClient *redis.Client, userRepo CacheUserRepository, fileService storage.FileService) (Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.BootstrapServers,
		"group.id":          config.GroupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		config:      config,
		consumer:    consumer,
		redisClient: redisClient,
		userRepo:    userRepo,
		fileService: fileService,
	}, nil
}

var logger = logging.GetLogger()

func (c *KafkaConsumer) Start() {
	err := c.consumer.SubscribeTopics(c.config.Topics, nil)
	if err != nil {
		logger.Errorf("Error subscribing to topics: %v", err)
		return
	}

	logger.Info("Kafka KafkaConsumer started...")

	for {
		msg, err := c.consumer.ReadMessage(-1)
		if err == nil {
			logger.Infof("Received message: %s", string(msg.Value))
			c.handleMessage(msg)
		} else {
			logger.Warnf("KafkaConsumer error: %v", err)
		}
	}
}

func (c *KafkaConsumer) Close() {
	c.consumer.Close()
}

func (c *KafkaConsumer) handleMessage(msg *kafka.Message) {
	var event events.Event
	err := json.Unmarshal(msg.Value, &event)
	if err != nil {
		logger.Errorf("Error parsing message: %v", err)
		return
	}

	switch event.Type {
	case "UserUpdated":
		break
	case "UserDeleted":
		break
	default:
		logger.Warnf("Unknown event type: %s", event.Type)
	}
}
