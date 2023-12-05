package main

import (
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"gopkg.in/robfig/cron.v2"
	"oursos.com/packages/alerts"
	"oursos.com/packages/api"
	"oursos.com/packages/db"
	image "oursos.com/packages/images"
	"oursos.com/packages/users"
	"oursos.com/packages/util"
)

func homeHandler(c echo.Context) error {
	return c.String(http.StatusOK, "OUR SOS BACKEND NOW MOVE ON")
}

func main() {
	err := godotenv.Load()
	util.CheckError(err)
	dbConn, err := db.Connection()
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // Allow all origins.
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	c := cron.New()
	c.AddFunc("@every 10m", func() {
		cleanAlertSQL := `DELETE FROM alerts WHERE time <= NOW() - INTERVAL '10 minute';`
		_, err = dbConn.Exec(cleanAlertSQL)
		util.CheckError(err)
	})

	c.Start()

	// db.SeedDatabase()
	e.GET("/", homeHandler)
	e.GET("/users", users.GetAllUsersHandler)
	e.GET("/users/:deviceId", users.GetUserById)
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
	e.POST("/news/:lang", api.GetNews)
	e.GET("/languages", api.ListSupportedLanguages)
	e.POST("/translateobject/:lang", api.TranslateObject)
	e.POST("/uploadimage", image.UploadImage)
	e.POST("/createuser", users.CreateUser)
	e.GET("/getuserbyid/:id", users.GetUserByUserId)
	e.GET("/getfriendsforuser/:id", users.GetFriendsForUsers)
	e.POST("/addfriend/:id1/:id2", users.AddFriend)
	e.Logger.Fatal(e.Start("0.0.0.0:" + os.Getenv("PORT")))
}
