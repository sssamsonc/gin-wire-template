package mysql_database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gin-wire-template/configs/common_config"
	"gin-wire-template/configs/mysql_db_config"
	"gin-wire-template/utils/log_util"
	"gin-wire-template/utils/once_util"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var (
	connector *Connector
	once      once_util.Once
	initError error
)

type Connector struct {
	client       *sql.DB
	dbConfig     mysql_db_config.Config
	commonConfig common_config.Config
	//errChan      <-chan error
}

func NewConnector(dbConfig mysql_db_config.Config, commonConfig common_config.Config) (*Connector, error) {
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
	sourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		c.dbConfig.Username,
		c.dbConfig.Password,
		c.dbConfig.Host,
		c.dbConfig.Database,
	)

	client, err := sql.Open("mysql", sourceName)
	if err != nil {
		log_util.Logger.Error("failed to connect to mysql::" + err.Error())
		return err
	}

	// Set database connection options
	client.SetMaxIdleConns(10)
	//client.SetMaxOpenConns(100)
	client.SetConnMaxLifetime(time.Minute * 3) //recommend timeout shorter than 5 minutes

	if err := c.pingClient(ctx, client); err != nil {
		return err
	}

	//errChan := make(chan error)
	//go c.monitor(ctx, client, errChan)
	c.client = client
	//c.errChan = errChan

	return nil
}

func (c *Connector) pingClient(ctx context.Context, client *sql.DB) error {
	if c.commonConfig.IsDemoMode { //ignore connection checking in demo mode
		return nil
	}

	if client != nil {
		if err := client.PingContext(ctx); err != nil {
			return errors.New("failed to ping mysql::" + err.Error())
		}
	}

	return nil
}

//// monitor the client connection per 60 seconds
//func (c *Connector) monitor(ctx context.Context, client *sql.DB, errChan chan error) {
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

func (c *Connector) Disconnect() error {
	if c.client == nil {
		return nil
	}

	return c.client.Close()
}

func (c *Connector) GetClient(ctx context.Context) (*sql.DB, error) {
	if c.client == nil {
		return nil, errors.New("mysql client is nil")
	}

	if err := c.pingClient(ctx, c.client); err != nil {
		return nil, err
	}

	return c.client, nil
}
