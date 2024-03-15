package mongo_database

import (
	"context"
	"errors"
	"fmt"
	"gin-wire-template/configs/common_config"
	"gin-wire-template/configs/mongo_db_config"
	"gin-wire-template/utils/log_util"
	"gin-wire-template/utils/once_util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	connector *Connector
	once      once_util.Once
	initError error
)

type Connector struct {
	client       *mongo.Client
	dbConfig     mongo_db_config.Config
	commonConfig common_config.Config
	//errChan      <-chan error
}

func NewConnector(dbConfig mongo_db_config.Config, commonConfig common_config.Config) (*Connector, error) {
	once.Do(func() {
		connector = &Connector{
			dbConfig:     dbConfig,
			commonConfig: commonConfig,
		}

		initError = connector.initializeClient(context.Background())
	})

	return connector, initError
}

func (c *Connector) initializeClient(ctx context.Context) error {
	uri := fmt.Sprintf("mongodb://%s:%s@%s/admin",
		c.dbConfig.Username,
		c.dbConfig.Password,
		c.dbConfig.Host,
	)
	if c.dbConfig.Replicaset != "" {
		uri += "?replicaSet=" + c.dbConfig.Replicaset
	}

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log_util.Logger.Error("failed to connect to mongo::" + err.Error())
		return err
	}

	if err := c.pingClient(ctx, client); err != nil {
		return err
	}

	//errChan := make(chan error)
	//go c.monitor(ctx, client, errChan)
	c.client = client
	//c.errChan = errChan

	return nil
}

func (c *Connector) pingClient(ctx context.Context, client *mongo.Client) error {
	if c.commonConfig.IsDemoMode { //ignore connection checking in demo mode
		return nil
	}

	if client != nil {
		if err := client.Ping(ctx, nil); err != nil {
			return errors.New("failed to ping mongo::" + err.Error())
		}
	}

	return nil
}

//// monitor the client connection per 60 seconds
//func (c *Connector) monitor(ctx context.Context, client *mongo.Client, errChan chan error) {
//	for {
//		if err := c.pingClient(ctx, client); err != nil {
//			errChan <- err
//			break
//		}
//		time.Sleep(60 * time.Second)
//	}
//}
//
//func (c *Connector) Err() <-chan error {
//	return c.errChan
//}

func (c *Connector) Disconnect(ctx context.Context) error {
	if c.client == nil {
		return nil
	}

	return c.client.Disconnect(ctx)
}

func (c *Connector) GetClient(ctx context.Context) (*mongo.Client, error) {
	if c.client == nil {
		return nil, errors.New("mongo client is nil")
	}

	if err := c.pingClient(ctx, c.client); err != nil {
		return nil, err
	}

	return c.client, nil
}

func (c *Connector) GetDb(ctx context.Context) (*mongo.Database, error) {
	if _, err := c.GetClient(ctx); err != nil {
		return nil, err
	}

	return c.client.Database(c.dbConfig.Database), nil
}
