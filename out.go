package main

import (
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"oursos.com/packages/alerts"
	"oursos.com/packages/api"
	"oursos.com/packages/db"
	"oursos.com/packages/users"
	"oursos.com/packages/util"
)

func homeHandler(c echo.Context) error {
	return c.String(http.StatusOK, "OUR SOS BACKEND NOW MOVE ON")
}

func main() {
	err := godotenv.Load()
	util.CheckError(err)
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // Allow all origins.
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	db.SeedDatabase()
	e.GET("/", homeHandler)
	e.GET("/users", users.GetAllUsersHandler)
	e.GET("/users/:id", users.GetUserById)
	e.DELETE("/deleteuser/:id", users.DeleteUser)
	e.PUT("/updateuser/:id", users.UpdateUser)
	e.GET("/alerts", alerts.GetAllAlertsHandler)
	e.GET("/alerts/:id", alerts.GetAlertById)
	e.POST("/reportalert", alerts.ReportAlert)
	e.GET("/aispeech", api.Speech)
	e.GET("/fires", api.GetForestFireData)
	e.GET("/earthquakes", api.GetEarthQuakes)
	e.POST("/translate", api.Translate)
	e.POST("/chat", api.ChatHandler)
	e.GET("/news", api.GetNews)
	e.GET("/languages", api.ListSupportedLanguages)
	e.POST("/translateobject", api.TranslateObject)

	e.Logger.Fatal(e.Start("0.0.0.0:" + os.Getenv("PORT")))
}
