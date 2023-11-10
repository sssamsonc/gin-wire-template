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

var connector *Connector

type Connector struct {
	client       *sql.DB
	error        error
	once         once_util.Once
	dbConfig     mysql_db_config.Config
	commonConfig common_config.Config
}

func NewConnector(dbConfig mysql_db_config.Config, commonConfig common_config.Config) (*Connector, error) {
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

func (c *Connector) verifyAndGetClient(ctx context.Context) (*sql.DB, error) {
	if c.commonConfig.IsDemoMode { //force to ignore the database connection
		return nil, nil
	}

	c.once.Do(func() {
		c.error = nil

		sourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s",
			c.dbConfig.Username,
			c.dbConfig.Password,
			c.dbConfig.Host,
			c.dbConfig.Database,
		)

		db, err := sql.Open("mysql", sourceName)
		if err != nil {
			c.error = err
			return
		}

		// Set database connection options
		db.SetMaxIdleConns(10)
		//db.SetMaxOpenConns(100)
		db.SetConnMaxLifetime(time.Minute * 3) //recommend timeout shorter than 5 minutes

		c.client = db
	})

	if c.client != nil {
		// Check the connection
		if err := c.client.PingContext(ctx); err != nil {
			c.error = err
			c.client = nil
		}
	}

	if c.error != nil {
		log_util.Logger.Error("fail to connect mysql database::" + c.error.Error())
		c.once.Reset()
	}

	return c.client, c.error
}

func (c *Connector) GetDb(ctx context.Context) (*sql.DB, error) {
	//check the client again, recreate if fail
	if _, err := connector.verifyAndGetClient(ctx); err != nil {
		return nil, err
	}

	if c.client == nil {
		return nil, errors.New("mysql client is nil")
	}

	return c.client, nil
}
