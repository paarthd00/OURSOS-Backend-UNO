package alerts

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"oursos.com/packages/db"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"oursos.com/packages/util"
)

func GetAlertById(c echo.Context) error {
	id := c.Param("id")

	db, err := db.Connection()
	util.CheckError(err)
	defer db.Close()

	query := "SELECT * FROM alerts WHERE id = $1"
	rows, err := db.Query(query, id)
	util.CheckError(err)
	defer rows.Close()

	if !rows.Next() {
		util.CheckError(err)
	}

	var alert Alert

	err = rows.Scan(&alert.ID, &alert.Message, &alert.Category, &alert.Severity, &alert.Time, &alert.Latitude, &alert.Longitude, &alert.Radius)

	util.CheckError(err)

	return c.JSON(http.StatusOK, alert)

}
