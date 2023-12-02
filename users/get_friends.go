package users

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
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

	user.Friends = friendIDs // Assign the parsed friend IDs to the 'Friends' field of the user

	// Query all friends at once
	friendsQuery := "SELECT * FROM users WHERE id = ANY($1)"
	friendsRows, err := db.Query(friendsQuery, pq.Array(friendIDs))
	util.CheckError(err)
	defer friendsRows.Close()

	// Scan each friend and append to a new slice
	var friends []User
	for friendsRows.Next() {
		var friend User
		var friendFriendsStr string
		errScan := friendsRows.Scan(&friend.ID, &friend.DeviceId, &friend.Username, &friend.Lat, &friend.Long, &friend.LanguagePreference, &friendFriendsStr, &friend.Profile)
		util.CheckError(errScan)

		// Parse the "friends" array from the string to []int for each friend
		friendFriendIDs, err := parseIntArray(friendFriendsStr)
		if err != nil {
			log.Fatal(err)
			return c.String(http.StatusInternalServerError, "Database error")
		}
		util.CheckError(err)

		friend.Friends = friendFriendIDs // Assign the parsed friend IDs to the 'Friends' field of each friend
		friends = append(friends, friend)
	}

	return c.JSON(http.StatusOK, friends)
}
