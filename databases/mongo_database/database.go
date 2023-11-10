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

var connector *Connector

type Connector struct {
	client       *mongo.Client
	error        error
	once         once_util.Once
	dbConfig     mongo_db_config.Config
	commonConfig common_config.Config
}

func NewConnector(dbConfig mongo_db_config.Config, commonConfig common_config.Config) (*Connector, error) {
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

func (c *Connector) verifyAndGetClient(ctx context.Context) (*mongo.Client, error) {
	if c.commonConfig.IsDemoMode { //force to ignore the database connection
		return nil, nil
	}

	c.once.Do(func() {
		c.error = nil

		// Set client options
		uri := fmt.Sprintf("mongodb://%s:%s@%s/admin",
			c.dbConfig.Username,
			c.dbConfig.Password,
			c.dbConfig.Host,
		)

		if c.dbConfig.Replicaset != "" {
			uri += "?replicaSet=" + c.dbConfig.Replicaset
		}

		// Apply to mongo URI
		clientOptions := options.Client().ApplyURI(uri)

		////set debug monitor
		//monitor := &event.CommandMonitor{
		//	Started: func(_ context.Context, e *event.CommandStartedEvent) {
		//		log_util.Logger.Info(string(e.Command))
		//	},
		//	Succeeded: func(_ context.Context, e *event.CommandSucceededEvent) {
		//		log_util.Logger.Info(string(e.Reply))
		//	},
		//	Failed: func(_ context.Context, e *event.CommandFailedEvent) {
		//		log_util.Logger.Info(e.Failure)
		//	},
		//}
		//
		//clientOptions.SetMonitor(monitor)

		// Connect
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			c.error = err
			return
		}

		c.client = client
	})

	if c.client != nil {
		// Check the connection
		if err := c.client.Ping(ctx, nil); err != nil {
			c.error = err
			c.client = nil
		}
	}

	if c.error != nil {
		log_util.Logger.Error("fail to connect mongo database::" + c.error.Error())
		c.once.Reset()
	}

	return c.client, c.error
}

func (c *Connector) GetDb(ctx context.Context) (*mongo.Database, error) {
	//check the client again, recreate if fail
	if _, err := connector.verifyAndGetClient(ctx); err != nil {
		return nil, err
	}

	if c.client == nil {
		return nil, errors.New("mongo client is nil")
	}

	return c.client.Database(c.dbConfig.Database), nil
}
