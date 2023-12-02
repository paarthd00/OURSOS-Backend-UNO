package users

import (
	"encoding/json"
	"strconv"

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

	// Start constructing the SQL query
	query := "UPDATE users SET "
	args := []interface{}{}
	i := 2

	// Iterate over the map and add each field to the query
	for field, value := range json_map {
		if field != "id" { // Exclude the 'id' field
			query += field + " = $" + strconv.Itoa(i) + ", "
			if field == "friends" {
				args = append(args, pq.Array(value.([]interface{})))
			} else {
				args = append(args, value)
			}
			i++
		}
	}

	// Remove the trailing comma and space, add the WHERE clause
	query = query[:len(query)-2] + " WHERE id = $1"
	args = append([]interface{}{id}, args...)

	stmt, err := dbConn.Prepare(query)
	util.CheckError(err)
	defer stmt.Close()

	_, err = stmt.Exec(args...)
	util.CheckError(err)

	return c.JSON(200, map[string]string{"message": "User updated successfully"})
}
