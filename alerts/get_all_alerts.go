package alerts

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"oursos.com/packages/db"

	"oursos.com/packages/util"
)

type Alert struct {
	ID       int     `json:"id"`
	Message  string  `json:"message"`
	Type     string  `json:"type"`
	Severity int8    `json:"severity"`
	Time     string  `json:"time"`
	Lat      float64 `json:"lat"`
	Long     float64 `json:"long"`
}

func GetAllAlertsHandler(c echo.Context) error {
	db, err := db.Connection()
	util.CheckError(err)
	rows, err := db.Query("SELECT * FROM alerts")
	util.CheckError(err)
	defer rows.Close()
	// client := redis.Client()
	// ctx := context.Background()

	// exists, err := client.Exists(ctx, "alerts").Result()
	// util.CheckError(err)
	var alerts []Alert

	// if exists == 1 {
	// alerts_json := client.Get(ctx, "alerts").Val()
	// err = json.Unmarshal([]byte(alerts_json), &alerts)
	// util.CheckError(err)
	// println("redis")
	// } else {
	for rows.Next() {
		var alert Alert
		if err := rows.Scan(&alert.ID, &alert.Message, &alert.Type, &alert.Severity, &alert.Time, &alert.Lat, &alert.Long); err != nil {
			util.CheckError(err)
		}
		alerts = append(alerts, alert)
	}
	// alerts_json, err := json.Marshal(alerts)
	util.CheckError(err)
	// rediserr := client.Set(ctx, "alerts", alerts_json, 0).Err()
	// util.CheckError(rediserr)
	// }

	return c.JSON(http.StatusOK, alerts)
}
