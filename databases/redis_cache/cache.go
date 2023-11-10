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
	"strconv"
	"time"
)

var connector *Connector

type Connector struct {
	client          *redis.Client
	error           error
	once            once_util.Once
	dbConfig        redis_config.Config
	commonConfig    common_config.Config
	isEnabled       bool
	isLockedToRetry bool
	failedAt        time.Time
}

func NewConnector(dbConfig redis_config.Config, commonConfig common_config.Config) (*Connector, error) {
	if connector == nil {
		connector = &Connector{
			dbConfig:     dbConfig,
			commonConfig: commonConfig,
		}
	}

	//init client
	_, err := connector.verifyAndGetClient(context.Background())
	return connector, err
}

func (c *Connector) verifyAndGetClient(ctx context.Context) (*redis.Client, error) {
	c.once.Do(func() {
		c.error = nil

		c.isEnabled = c.dbConfig.Host != ""

		if !c.isEnabled {
			return
		}

		if c.isLockedToRetry {
			return
		}

		dbNum, _ := strconv.Atoi(c.dbConfig.Database)

		opts := &redis.Options{
			Addr:     c.dbConfig.Host,
			Username: c.dbConfig.Username,
			Password: c.dbConfig.Password,
			DB:       dbNum,
		}

		c.client = redis.NewClient(opts)
	})

	if !c.isEnabled {
		log_util.Logger.Debug("redis cache is disabled")
		return nil, nil
	}

	if c.client != nil {
		// Check the connection
		if err := c.client.Ping(ctx).Err(); err != nil {
			c.error = err
			c.client = nil
		}
	}

	if c.error != nil {
		if c.isLockedToRetry {
			if time.Now().Sub(c.failedAt) > c.dbConfig.CacheRetryLockTime {
				c.isLockedToRetry = false
				c.once.Reset()
			} else {
				x := c.dbConfig.CacheRetryLockTime - time.Now().Sub(c.failedAt)
				log_util.Logger.Debug("redis cache will retry to connect in " + fmt.Sprintf("%02dm:%02ds", int(x.Minutes()), int(x.Seconds())%60))
			}
		} else {
			log_util.Logger.Debug("fail to connect redis cache (will retry to connect after " + fmt.Sprintf("%02dm:%02ds", int(c.dbConfig.CacheRetryLockTime.Minutes()), int(c.dbConfig.CacheRetryLockTime.Seconds())%60) + ") ::" + c.error.Error())
			c.failedAt = time.Now()
			c.isLockedToRetry = true
		}
	}

	return c.client, c.error
}

func (c *Connector) GetDb(ctx context.Context) (*redis.Client, error) {
	//check the client again, recreate if fail
	if _, err := connector.verifyAndGetClient(ctx); err != nil {
		return nil, err
	}

	if c.client == nil {
		return nil, errors.New("redis client is nil")
	}

	return c.client, nil
}
