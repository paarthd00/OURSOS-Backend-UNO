package users

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"oursos.com/packages/db"
	"oursos.com/packages/util"
)

func AddFriend(c echo.Context) error {
	// Get the user IDs from the URL parameters
	id1, err := strconv.Atoi(c.Param("id1"))
	util.CheckError(err)

	id2, err := strconv.Atoi(c.Param("id2"))
	util.CheckError(err)

	// Connect to the database
	db, err := db.Connection()
	util.CheckError(err)
	defer db.Close()

	// Fetch the users from the database
	user1, err := getUserById(db, id1)
	util.CheckError(err)

	user2, err := getUserById(db, id2)
	util.CheckError(err)

	// Add each user to the other's friends array
	user1.Friends = append(user1.Friends, id2)
	user2.Friends = append(user2.Friends, id1)

	// Update the users in the database
	err = updateUser(db, user1)
	util.CheckError(err)

	err = updateUser(db, user2)
	util.CheckError(err)

	// Return a success response
	return c.JSON(http.StatusOK, map[string]string{"message": "Successfully added friend"})
}

// You'll need to implement these functions to fetch and update users in the database
func getUserById(db *sql.DB, id int) (User, error) {
	var user User
	var friends sql.NullString
	err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&user.ID, &user.DeviceId, &user.Username, &user.Lat, &user.Long, &user.LanguagePreference, &friends, &user.Profile)
	if err != nil {
		return User{}, err
	}
	if friends.Valid {
		user.Friends = parseStringToIntSlice(friends.String)
	}
	return user, nil
}

func updateUser(db *sql.DB, user User) error {
	_, err := db.Exec("UPDATE users SET deviceId = $1, username = $2, lat = $3, long = $4, languagepreference = $5, friends = $6, profile = $7 WHERE id = $8", user.DeviceId, user.Username, user.Lat, user.Long, user.LanguagePreference, pq.Array(user.Friends), user.Profile, user.ID)
	return err
}

func parseStringToIntSlice(s string) []int {
	s = strings.Trim(s, "{}")
	if s == "" {
		return []int{}
	}
	ss := strings.Split(s, ",")
	result := make([]int, len(ss))
	for i, v := range ss {
		result[i], _ = strconv.Atoi(v)
	}
	return result
}
