package users

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"oursos.com/packages/db"
	"oursos.com/packages/util"
)

func DeleteUser(c echo.Context) error {

	db, err := db.Connection()
	util.CheckError(err)

	id := c.Param("id")
	defer db.Close()
	// Prepare and execute the DELETE query
	query := "DELETE FROM users WHERE id = $1"
	result, err := db.Exec(query, id)
	util.CheckError(err)

	// Check the number of rows affected
	rowsAffected, err := result.RowsAffected()
	util.CheckError(err)

	// Check if any rows were deleted
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
}
