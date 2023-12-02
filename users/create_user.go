package users

import (
	"encoding/json"
	"net/http"

	"github.com/lib/pq"

	"github.com/labstack/echo/v4"
	"oursos.com/packages/db"
	"oursos.com/packages/util"
)

func CreateUser(c echo.Context) error {
	dbConn, err := db.Connection()
	util.CheckError(err)
	defer dbConn.Close()

	json_map := make(map[string]interface{})
	errEnc := json.NewDecoder(c.Request().Body).Decode(&json_map)
	util.CheckError(errEnc)

	username := json_map["username"].(string)

	// Check if a user with the same username already exists
	userExistsQuery := "SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)"
	row := dbConn.QueryRow(userExistsQuery, username)
	var exists bool
	err = row.Scan(&exists)
	util.CheckError(err)

	if exists {
		return c.JSON(http.StatusConflict, map[string]string{"message": "Username already exists"})
	}

	deviceId := json_map["deviceId"].(string)
	lat := json_map["lat"].(float64)
	long := json_map["long"].(float64)
	languagepreference := json_map["languagepreference"].(string)
	profile := json_map["profile"].(string)

	// Check if friends exists in json_map, if not, set to empty array
	var friendsArr []int
	if friends, ok := json_map["friends"]; ok {
		friendsArr = make([]int, len(friends.([]interface{})))
		for i, v := range friends.([]interface{}) {
			friendsArr[i] = int(v.(float64))
		}
	}

	// Create a prepared statement
	stmt, err := dbConn.Prepare("INSERT INTO users (deviceId, username, lat, long, languagepreference, friends, profile) VALUES ($1, $2, $3, $4, $5, $6, $7)")
	util.CheckError(err)
	defer stmt.Close()

	_, err = stmt.Exec(deviceId, username, lat, long, languagepreference, pq.Array(friendsArr), profile)
	util.CheckError(err)

	return c.JSON(http.StatusOK, map[string]string{"message": "User created successfully"})
}
