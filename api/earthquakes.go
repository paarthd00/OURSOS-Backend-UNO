package api

import (
	"context"
	"net/http"
	"time"

	"encoding/json"

	"github.com/labstack/echo/v4"
	"oursos.com/packages/redis"
	"oursos.com/packages/util"
)

func GetEarthQuakes(c echo.Context) error {

	client := redis.Client()
	ctx := context.Background()

	exists, err := client.Exists(ctx, "earthquakes").Result()
	util.CheckError(err)
	var earthquakeJSON interface{}

	if exists == 1 {
		earthquakes_json := client.Get(ctx, "earthquakes").Val()
		err = json.Unmarshal([]byte(earthquakes_json), &earthquakeJSON)
		util.CheckError(err)
		println("redis")
	} else {
		currentTime := time.Now().UTC()

		yesterday := currentTime.AddDate(0, 0, -1)

		yesterdayDate := yesterday.Format("2006-01-02")
		todaysDate := currentTime.Format("2006-01-02")

		apiURL := "https://earthquake.usgs.gov/fdsnws/event/1/query?format=geojson&starttime=" + yesterdayDate + "&endtime=" + todaysDate

		req, _ := http.NewRequest("GET", apiURL, nil)

		res, err := http.DefaultClient.Do(req)
		util.CheckError(err)
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return c.JSON(res.StatusCode, map[string]string{"error": "Failed to retrieve earthquake data"})
		}
		// Read the response body
		json.NewDecoder(res.Body).Decode(&earthquakeJSON)

		util.CheckError(err)
		earthquakes_json, err := json.Marshal(earthquakeJSON)
		util.CheckError(err)
		rediserr := client.Set(ctx, "earthquakes", earthquakes_json, 0).Err()
		util.CheckError(rediserr)
	}

	return c.JSON(http.StatusOK, earthquakeJSON)
}
