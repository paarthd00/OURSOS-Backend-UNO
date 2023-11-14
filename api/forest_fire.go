package api

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"time"

	"github.com/labstack/echo/v4"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"oursos.com/packages/redis"
	"oursos.com/packages/util"
)

// FireData represents the structure of your CSV data
type FireData struct {
	CountryID  string `json:"country_id"`
	Latitude   string `json:"latitude"`
	Longitude  string `json:"longitude"`
	BrightTI4  string `json:"bright_ti4"`
	Scan       string `json:"scan"`
	Track      string `json:"track"`
	AcqDate    string `json:"acq_date"`
	AcqTime    string `json:"acq_time"`
	Satellite  string `json:"satellite"`
	Instrument string `json:"instrument"`
	Confidence string `json:"confidence"`
	Version    string `json:"version"`
	BrightTI5  string `json:"bright_ti5"`
	FRP        string `json:"frp"`
	DayNight   string `json:"daynight"`
}

func GetForestFireData(c echo.Context) error {
	err := godotenv.Load()

	util.CheckError(err)
	currentTime := time.Now().UTC()
	client := redis.Client()
	ctx := context.Background()

	exists, err := client.Exists(ctx, "fires").Result()
	util.CheckError(err)

	var fireDataList []FireData

	if exists == 1 {
		fireDataList_json := client.Get(ctx, "fires").Val()
		err = json.Unmarshal([]byte(fireDataList_json), &fireDataList)
		util.CheckError(err)
		println("redis")
	} else {

		todaysDate := currentTime.AddDate(0, 0, -1).Format("2006-01-02")
		key := os.Getenv("NASA_FIRE_KEY")
		apiURL := "https://firms.modaps.eosdis.nasa.gov/api/country/csv/" + key + "/VIIRS_SNPP_NRT/CAN/1/" + todaysDate
		response, err := http.Get(apiURL)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		// Create a CSV reader to read the response body
		reader := csv.NewReader(response.Body)

		// Read and parse the CSV records
		var records [][]string
		for {
			record, err := reader.Read()
			if err != nil {
				break
			}
			// Replace stray double-quote characters in each field
			for i := range record {
				record[i] = strings.Replace(record[i], `"`, ``, -1)
			}
			records = append(records, record)
		}

		// Iterate through the CSV records and populate the FireData struct
		for _, record := range records[1:] { // Skip the header row
			fireData := FireData{
				CountryID:  record[0],
				Latitude:   record[1],
				Longitude:  record[2],
				BrightTI4:  record[3],
				Scan:       record[4],
				Track:      record[5],
				AcqDate:    record[6],
				AcqTime:    record[7],
				Satellite:  record[8],
				Instrument: record[9],
				Confidence: record[10],
				Version:    record[11],
				BrightTI5:  record[12],
				FRP:        record[13],
				DayNight:   record[14],
			}
			fireDataList = append(fireDataList, fireData)
		}
		fireDataList_json, err := json.Marshal(fireDataList)
		util.CheckError(err)
		rediserr := client.Set(ctx, "fires", fireDataList_json, 0).Err()
		util.CheckError(rediserr)
	}

	return c.JSON(http.StatusOK, fireDataList)
}
