package users

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"oursos.com/packages/db"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"oursos.com/packages/util"
)

func GetUserByUserId(c echo.Context) error {
	id := c.Param("id")

	db, err := db.Connection()
	util.CheckError(err)
	defer db.Close()

	query := "SELECT * FROM users WHERE id = $1"
	rows, err := db.Query(query, id)

	util.CheckError(err)
	defer rows.Close()

	if !rows.Next() {
		return c.JSON(http.StatusOK, nil)
	}

	var user User
	var friendsStr sql.NullString // Use sql.NullString instead of string
	errScan := rows.Scan(&user.ID, &user.DeviceId, &user.Username, &user.Lat, &user.Long, &user.LanguagePreference, &friendsStr, &user.Profile)
	util.CheckError(errScan)

	// Check if friendsStr is not NULL before parsing
	if friendsStr.Valid {
		// Parse the "friends" array from the string to []int
		user.Friends, err = parseIntArray(friendsStr.String)
		if err != nil {
			log.Fatal(err)
			return c.String(http.StatusInternalServerError, "Database error")
		}
		util.CheckError(err)
	}

	return c.JSON(http.StatusOK, user)
}
