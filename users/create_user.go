package users

import (
	"encoding/json"

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

	deviceId := json_map["deviceId"]
	username := json_map["username"]
	lat := json_map["lat"]
	long := json_map["long"]
	languagepreference := json_map["languagepreference"]
	friends := json_map["friends"]
	profile := json_map["profile"]

	// Convert friends slice to []int
	friendsArr := make([]int, len(friends.([]interface{})))
	for i, v := range friends.([]interface{}) {
		friendsArr[i] = int(v.(float64))
	}

	// Create a prepared statement
	stmt, err := dbConn.Prepare("INSERT INTO users (deviceId, username, lat, long, languagepreference, friends, profile) VALUES ($1, $2, $3, $4, $5, $6, $7)")
	util.CheckError(err)
	defer stmt.Close()

	_, err = stmt.Exec(deviceId, username, lat, long, languagepreference, pq.Array(friendsArr), profile)
	util.CheckError(err)

	return c.JSON(200, map[string]string{"message": "User created successfully"})
}
