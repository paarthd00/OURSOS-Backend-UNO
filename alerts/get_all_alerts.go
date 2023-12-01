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

	// Check if the alerts table exists
	tableCheck, err := db.Query("SELECT 1 FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'alerts'")
	util.CheckError(err)
	if !tableCheck.Next() {
		return c.JSON(http.StatusOK, make([]Alert, 0))
	}

	rows, err := db.Query("SELECT * FROM alerts")
	util.CheckError(err)
	defer rows.Close()

	var alerts []Alert

	for rows.Next() {
		var alert Alert
		if err := rows.Scan(&alert.ID, &alert.Message, &alert.Type, &alert.Severity, &alert.Time, &alert.Lat, &alert.Long); err != nil {
			util.CheckError(err)
		}
		alerts = append(alerts, alert)
	}

	if err = rows.Err(); err != nil {
		util.CheckError(err)
	}

	if len(alerts) == 0 {
		alerts = make([]Alert, 0)
	}

	return c.JSON(http.StatusOK, alerts)
}
