package text_menu_repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"gin-wire-template/models/text_menu"
	"gin-wire-template/utils/http_util"
	"gin-wire-template/utils/log_util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"strings"
)

// GetMenu
// TODO - Example for getting data from MongoDB with Redis cache
func (repo *Repository) GetMenu(ctx context.Context, itemTypes []int) ([]text_menu.TextMenu, error) {
	var results []text_menu.TextMenu
	rdb, rdbErr := repo.redisConnector.GetDb(ctx)

	if rdbErr == nil {
		data, _ := rdb.Get(ctx, TEXT_MENU_CACHE_NAME).Bytes()
		if err := json.Unmarshal(data, &results); err == nil {
			log_util.Logger.Debug("Used Cache! " + TEXT_MENU_CACHE_NAME)

			err := rdb.Expire(ctx, TEXT_MENU_CACHE_NAME, repo.redisConfig.CacheTTL).Err()
			if err == nil {
				log_util.Logger.Debug("Renewed cache TTL! " + TEXT_MENU_CACHE_NAME)
			}

			return results, nil
		}
	}

	db, err := repo.mongoConnector.GetDb(ctx)
	if err != nil {
		return nil, err
	}

	collection := db.Collection(COLLECTION_TEXT_MENU)
	filter := bson.M{}
	if len(itemTypes) > 0 {
		filter["type"] = bson.M{"$in": itemTypes}
	}

	opts := options.Find().SetSort(bson.D{{"order_id", 1}})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if rdbErr == nil && len(results) > 0 {
		if res, err := json.Marshal(results); err == nil {
			//set cache
			exist, _ := rdb.Exists(ctx, TEXT_MENU_CACHE_NAME).Result()
			if exist == 0 {
				rdb.Set(ctx, TEXT_MENU_CACHE_NAME, res, repo.redisConfig.CacheTTL)
				log_util.Logger.Debug("Saved to Cache! " + TEXT_MENU_CACHE_NAME)
			}
		}
	}

	return results, nil
}

// GetMenuFromMySQL
// TODO - Example for getting data from MySQL
func (repo *Repository) GetMenuFromMySQL(ctx context.Context, itemTypes []int) ([]text_menu.TextMenu, error) {
	var (
		results []text_menu.TextMenu
		rows    *sql.Rows
	)

	db, err := repo.mysqlConnector.GetDb(ctx)
	if err != nil {
		return nil, err
	}

	if len(itemTypes) > 0 {
		//example - prepared statement

		//handle multiple items
		query := fmt.Sprintf(
			"SELECT * FROM "+COLLECTION_TEXT_MENU+" WHERE type IN (?) ORDER BY order_id",
			strings.Repeat("?,", len(itemTypes)-1)+"?")

		stmt, err := db.Prepare(query)
		if err != nil {
			return nil, err
		}

		//change the item type to interface
		genericItemType := make([]interface{}, len(itemTypes))
		for i, v := range itemTypes {
			genericItemType[i] = v
		}

		rows, err = stmt.Query(genericItemType...)
		if err != nil {
			return nil, err
		}

	} else {
		//example - raw statement
		query := "SELECT * FROM " + COLLECTION_TEXT_MENU + " ORDER BY order_id"
		rows, err = db.Query(query)
		if err != nil {
			return nil, err
		}
	}

	defer rows.Close()

	for rows.Next() {
		var r text_menu.TextMenu

		err = rows.Scan(
			&r.Id,
			&r.Type,
			&r.Name,
			&r.Url,
		)
		if err != nil {
			log_util.Logger.Error("Error occurred when scan rows::" + err.Error())
		} else {
			results = append(results, r)
		}
	}
	return results, nil
}

// GetMenuFromExternalAPI
// TODO - Example for getting data from external api
func (repo *Repository) GetMenuFromExternalAPI(ctx context.Context) ([]text_menu.TextMenu, error) {

	response, err := http_util.NewClient().Get("EXTERNAL_API_URL")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("api response is not ok: " + response.Status)
	}

	var results []text_menu.TextMenu

	if err := json.NewDecoder(response.Body).Decode(&results); err != nil {
		return nil, err
	}

	return results, err
}
