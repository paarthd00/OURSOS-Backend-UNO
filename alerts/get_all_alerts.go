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
	ID        int    `json:"id"`
	Message   string `json:"message"`
	Category  string `json:"category"`
	Severity  string `json:"severity"`
	Time      string `json:"time"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Radius    string `json:"radius"`
}

func GetAllAlertsHandler(c echo.Context) error {
	db, err := db.Connection()
	util.CheckError(err)
	rows, err := db.Query("SELECT * FROM alerts")
	util.CheckError(err)
	defer rows.Close()

	var alerts []Alert

	for rows.Next() {
		var alert Alert
		if err := rows.Scan(&alert.ID, &alert.Message, &alert.Category, &alert.Severity, &alert.Time, &alert.Latitude, &alert.Longitude, &alert.Radius); err != nil {
			util.CheckError(err)
		}
		alerts = append(alerts, alert)
	}

	util.CheckError(err)

	return c.JSON(http.StatusOK, alerts)
}
