package api

import (
	"net/http"
	"time"

	"encoding/json"

	"github.com/labstack/echo/v4"
	"oursos.com/packages/util"
)

func GetEarthQuakes(c echo.Context) error {
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
	var bodyJSON interface{}
	// Read the response body
	json.NewDecoder(res.Body).Decode(&bodyJSON)

	util.CheckError(err)
	return c.JSON(http.StatusOK, bodyJSON)
}
