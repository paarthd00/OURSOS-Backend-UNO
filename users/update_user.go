package users

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"oursos.com/packages/db"
	"oursos.com/packages/util"
)

func UpdateUser(c echo.Context) error {

	dbConn, err := db.Connection()
	util.CheckError(err)
	defer dbConn.Close()

	id := c.Param("id")
	json_map := make(map[string]interface{})
	errEnc := json.NewDecoder(c.Request().Body).Decode(&json_map)
	util.CheckError(errEnc)

	username := json_map["username"]
	lat := json_map["lat"]
	long := json_map["long"]
	languagepreference := json_map["languagepreference"]
	friends := json_map["friends"].([]interface{})

	stmt, err := dbConn.Prepare("UPDATE users SET username = $2, lat = $3, long=$4, languagepreference = $5, friends = $6 WHERE id = $1 ")
	util.CheckError(err)
	defer stmt.Close()

	_, err = stmt.Exec(id, username, lat, long, languagepreference.(string), pq.Array(friends))
	util.CheckError(err)

	return c.JSON(200, map[string]string{"message": "User updated successfully"})
}
