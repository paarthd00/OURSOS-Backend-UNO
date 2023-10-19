package users

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"oursos.com/packages/db"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"oursos.com/packages/util"
)

func GetUserById(c echo.Context) error {

	id := c.Param("id")

	db, err := db.Connection()
	util.CheckError(err)
	defer db.Close()

	query := "SELECT * FROM users WHERE id = $1"
	rows, err := db.Query(query, id)

	util.CheckError(err)
	defer rows.Close()

	if !rows.Next() {
		util.CheckError(fmt.Errorf("no user with id %s", id))
	}

	var user User
	var friendsStr string
	errScan := rows.Scan(&user.ID, &user.Username, pq.Array(&user.Locations), &user.LanguagePreference, &friendsStr)
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
