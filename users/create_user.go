package users

import (
	"encoding/json"

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
	long := json_map["lang"]
	languagepreference := json_map["languagepreference"]
	friends := json_map["friends"]

	// Create a prepared statement
	stmt, err := dbConn.Prepare("INSERT INTO users (deviceId, username, lat, long, languagepreference, friends) VALUES (?, ?, ?, ?)")
	util.CheckError(err)
	defer stmt.Close()

	_, err = stmt.Exec(deviceId, username, lat, long, languagepreference, friends)
	util.CheckError(err)

	return c.JSON(200, map[string]string{"message": "User created successfully"})

}
