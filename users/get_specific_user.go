package users

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"oursos.com/packages/db"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"oursos.com/packages/util"
)

func GetUserById(c echo.Context) error {

	deviceId := c.Param("deviceId")

	db, err := db.Connection()
	util.CheckError(err)
	defer db.Close()

	query := "SELECT * FROM users WHERE deviceId = $1"
	rows, err := db.Query(query, deviceId)

	util.CheckError(err)
	defer rows.Close()

	if !rows.Next() {
		return c.JSON(http.StatusOK, nil)
	}

	var user User
	var friendsStr string
	errScan := rows.Scan(&user.ID, &user.DeviceId, &user.Username, &user.Lat, &user.Long, &user.LanguagePreference, &friendsStr, &user.Profile)
	util.CheckError(errScan)

	// Parse the "friends" array from the string to []int
	user.Friends, err = parseIntArray(friendsStr)
	if err != nil {
		log.Fatal(err)
		return c.String(http.StatusInternalServerError, "Database error")
	}

	util.CheckError(err)

	return c.JSON(http.StatusOK, user)

}
