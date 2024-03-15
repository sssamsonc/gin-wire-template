package redis_cache

import (
	"context"
	"errors"
	"fmt"
	"gin-wire-template/configs/common_config"
	"gin-wire-template/configs/redis_config"
	"gin-wire-template/utils/log_util"
	"gin-wire-template/utils/once_util"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	connector *Connector
	once      once_util.Once
)

type Connector struct {
	client       *redis.Client
	dbConfig     redis_config.Config
	commonConfig common_config.Config
	isEnabled    bool
	isLocked     bool
	failedAt     time.Time
}

func NewConnector(dbConfig redis_config.Config, commonConfig common_config.Config) *Connector {
	once.Do(func() {
		connector = &Connector{
			dbConfig:     dbConfig,
			commonConfig: commonConfig,
			isEnabled:    dbConfig.Host != "",
		}

		connector.initializeClient()
	})

	return connector
}

func (c *Connector) initializeClient() {
	if !c.isEnabled {
		log_util.Logger.Debug("redis is disabled in configuration")
		return
	}

	opts := &redis.Options{
		Addr:     c.dbConfig.Host,
		Username: c.dbConfig.Username,
		Password: c.dbConfig.Password,
		DB:       c.dbConfig.Database,
	}

	c.client = redis.NewClient(opts)
}

func (c *Connector) lockDown() {
	if !c.isLocked {
		if c.client != nil {
			c.client.Close()
		}
		c.client = nil
		c.isLocked = true
		c.failedAt = time.Now()
		log_util.Logger.Warn("redis will be available to retry after " + fmt.Sprintf("%02dm:%02ds", int(c.dbConfig.CacheRetryLockTime.Minutes()), int(c.dbConfig.CacheRetryLockTime.Seconds())%60))
	}
}

func (c *Connector) GetClient(ctx context.Context) (*redis.Client, error) {
	if !c.isEnabled {
		return nil, errors.New("redis is disabled")
	}

	if c.isLocked {
		if time.Now().Sub(c.failedAt) > c.dbConfig.CacheRetryLockTime {
			c.isLocked = false
			c.initializeClient() // Attempt re-initialization if needed
		} else {
			x := c.dbConfig.CacheRetryLockTime - time.Now().Sub(c.failedAt)
			log_util.Logger.Debug("redis is temporarily locked due to previous errors ||remain::" + fmt.Sprintf("%02dm:%02ds", int(x.Minutes()), int(x.Seconds())%60))
		}
	}

	if c.client != nil {
		if err := c.client.Ping(ctx).Err(); err != nil {
			log_util.Logger.Warn("redis ping failed::" + err.Error())
			c.lockDown()
			return nil, errors.New("failed to connect to redis")
		}
	} else {
		return nil, errors.New("redis client is nil")
	}

	return c.client, nil
}
