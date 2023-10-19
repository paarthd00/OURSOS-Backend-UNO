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

	username := json_map["username"]
	locations := json_map["locations"]
	languagepreference := json_map["languagepreference"]
	friends := json_map["friends"]

	// Create a prepared statement
	stmt, err := dbConn.Prepare("INSERT INTO users (username, locations, languagepreference, friends) VALUES (?, ?, ?, ?)")
	util.CheckError(err)
	defer stmt.Close()

	_, err = stmt.Exec(username, locations, languagepreference, friends)
	util.CheckError(err)

	return c.JSON(200, map[string]string{"message": "User created successfully"})

}
