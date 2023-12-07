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

	jsonMap := make(map[string]interface{})
	errEnc := json.NewDecoder(c.Request().Body).Decode(&jsonMap)
	util.CheckError(errEnc)

	username := jsonMap["username"].(string)

	// Check if a user with the same username already exists
	userExistsQuery := "SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)"
	row := dbConn.QueryRow(userExistsQuery, username)
	var exists bool
	err = row.Scan(&exists)
	util.CheckError(err)

	if exists {
		return c.JSON(http.StatusConflict, map[string]string{"message": "Username already exists"})
	}

	deviceID := jsonMap["deviceId"].(string)
	lat := jsonMap["lat"].(float64)
	long := jsonMap["long"].(float64)
	languagePreference := jsonMap["languagepreference"].(string)
	profile := jsonMap["profile"].(string)

	// Check if friends exist in jsonMap, if not, set to empty array
	var friendsArr []int
	if friends, ok := jsonMap["friends"]; ok {
		friendsArr = make([]int, len(friends.([]interface{})))
		for i, v := range friends.([]interface{}) {
			friendsArr[i] = int(v.(float64))
		}
	}

	// Create a prepared statement
	stmt, err := dbConn.Prepare("INSERT INTO users (deviceId, username, lat, long, languagepreference, friends, profile) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id")
	util.CheckError(err)
	defer stmt.Close()

	var userID int
	err = stmt.QueryRow(deviceID, username, lat, long, languagePreference, pq.Array(friendsArr), profile).Scan(&userID)
	util.CheckError(err)

	createdUser := User{
		ID:                 userID,
		DeviceId:           deviceID,
		Username:           username,
		Lat:                lat,
		Long:               long,
		LanguagePreference: languagePreference,
		Friends:            friendsArr,
		Profile:            profile,
	}

	return c.JSON(http.StatusOK, createdUser)
}
