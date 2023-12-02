package users

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"oursos.com/packages/db"
	"oursos.com/packages/util"
)

func GetFriendsForUsers(c echo.Context) error {
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
	var friendsStr string
	errScan := rows.Scan(&user.ID, &user.DeviceId, &user.Username, &user.Lat, &user.Long, &user.LanguagePreference, &friendsStr, &user.Profile)
	util.CheckError(errScan)

	// Parse the "friends" array from the string to []int
	friendIDs, err := parseIntArray(friendsStr)
	if err != nil {
		log.Fatal(err)
		return c.String(http.StatusInternalServerError, "Database error")
	}
	util.CheckError(err)

	// Query each friend ID and append to a new slice
	var friends []User
	for _, friendID := range friendIDs {
		friendQuery := "SELECT * FROM users WHERE id = $1"
		friendRows, err := db.Query(friendQuery, friendID)
		util.CheckError(err)

		if friendRows.Next() {
			var friend User
			errScan := friendRows.Scan(&friend.ID, &friend.DeviceId, &friend.Username, &friend.Lat, &friend.Long, &friend.LanguagePreference, &friendsStr, &friend.Profile)
			util.CheckError(errScan)
			friends = append(friends, friend)
		}
		friendRows.Close()
	}

	return c.JSON(http.StatusOK, friends)
}
