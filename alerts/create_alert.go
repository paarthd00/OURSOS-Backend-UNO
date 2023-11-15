package alerts

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
	"oursos.com/packages/db"
	"oursos.com/packages/util"
)

func ReportAlert(c echo.Context) error {

	dbConn, err := db.Connection()
	util.CheckError(err)
	defer dbConn.Close()

	var alert Alert
	err = json.NewDecoder(c.Request().Body).Decode(&alert)
	util.CheckError(err)

	insertAlertSQL := `
        INSERT INTO alerts (message, category, severity, latitude, longitude, radius)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err = dbConn.Exec(insertAlertSQL, alert.Message, alert.Category, alert.Severity, alert.Latitude, alert.Longitude, alert.Radius)

	util.CheckError(err)
	return c.JSON(200, map[string]string{"message": "Alert created successfully"})
}
