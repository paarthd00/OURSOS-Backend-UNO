package users

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"oursos.com/packages/db"
	"oursos.com/packages/util"
)

func RemoveFriend(c echo.Context) error {
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

	// Check if they are friends
	if !areFriends(user1.Friends, id2) || !areFriends(user2.Friends, id1) {
		// Return a response indicating they are not friends
		return c.JSON(http.StatusOK, map[string]string{"message": "Users are not friends"})
	}

	// Remove friend ID from each user's Friends array
	user1.Friends = removeFriend(user1.Friends, id2)
	user2.Friends = removeFriend(user2.Friends, id1)

	// Update the users in the database
	err = updateUser(db, user1)
	util.CheckError(err)

	err = updateUser(db, user2)
	util.CheckError(err)

	// Return a success response
	return c.JSON(http.StatusOK, map[string]string{"message": "Successfully removed friend"})
}

// Function to remove a friend ID from the Friends array
func removeFriend(friends []int, friendID int) []int {
	var updatedFriends []int
	for _, id := range friends {
		if id != friendID {
			updatedFriends = append(updatedFriends, id)
		}
	}
	return updatedFriends
}
