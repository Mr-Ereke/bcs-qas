package redis

import (
	"context"
	"net"
	"strconv"
	"sync"
	"time"

	"app/models"

	"github.com/redis/go-redis/v9"
	"gitlab.online-fx.com/go-packages/logger"
)

type Client struct {
	mutex  sync.RWMutex
	ctx    context.Context
	config models.RedisConfig
	client *redis.Client
}

func New(config *models.Config) *Client {
	redisClient := &Client{
		ctx:    context.Background(),
		config: config.Redis,
	}

	redisClient.CreateClient()
	go redisClient.ping()

	return redisClient
}

func (c *Client) CreateClient() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	address := net.JoinHostPort(c.config.Address, strconv.Itoa(c.config.Port))
	c.client = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: c.config.Password,
		DB:       c.config.DB,
	})
	logger.Info("New Redis client created")
}

func (c *Client) ping() {
	for range time.Tick(time.Second * time.Duration(c.config.PingTimeout)) {
		err := c.client.Ping(c.ctx).Err()
		if err != nil {
			logger.Errorf("Failed Redis ping. Error: %v", err)
			time.Sleep(time.Second) // delay for recreate client
			c.CreateClient()
		}
	}
}
