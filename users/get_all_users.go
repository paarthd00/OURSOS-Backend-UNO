package users

import (
	"net/http"
	"strconv"
	"strings"

	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"oursos.com/packages/db"
	"oursos.com/packages/util"
)

type User struct {
	ID                 int     `json:"id"`
	DeviceId           string  `json:"deviceId"`
	Username           string  `json:"username"`
	Lat                float64 `json:"lat"`
	Long               float64 `json:"long"`
	LanguagePreference string  `json:"languagepreference"`
	Friends            []int   `json:"friends"`
	Profile            string  `json:"profile"`
}

func GetAllUsersHandler(c echo.Context) error {

	db, err := db.Connection()
	util.CheckError(err)

	rows, err := db.Query("SELECT id, deviceId, username, lat, long, languagepreference, friends, profile FROM users")
	if err != nil {
		log.Fatal(err)
		return c.String(http.StatusInternalServerError, "Database error")
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		var friendsStr string // To hold the array as a string
		err := rows.Scan(&user.ID, &user.DeviceId, &user.Username, &user.Lat, &user.Long, &user.LanguagePreference, &friendsStr, &user.Profile)
		if err != nil {
			log.Fatal(err)
			return c.String(http.StatusInternalServerError, "Database error")
		}

		// Parse the "friends" array from the string to []int
		user.Friends, err = parseIntArray(friendsStr)
		if err != nil {
			log.Fatal(err)
			return c.String(http.StatusInternalServerError, "Database error")
		}

		users = append(users, user)
	}

	return c.JSON(http.StatusOK, users)

}

func parseIntArray(input string) ([]int, error) {
	// Split the input string by commas
	parts := strings.Split(input, ",")

	var result []int
	for _, part := range parts {
		// Parse each part as an integer
		value, err := strconv.Atoi(strings.Trim(part, "{}"))
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}

	return result, nil
}
