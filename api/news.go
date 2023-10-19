package api

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	g "github.com/serpapi/google-search-results-golang"
	"oursos.com/packages/util"
)

type NewsArticle struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	PublishedAt string `json:"publishedAt"`
	Source      struct {
		Name string `json:"name"`
	} `json:"source"`
}

type NewsResponse struct {
	Articles []NewsArticle `json:"articles"`
}

func GetNews(c echo.Context) error {
	err := godotenv.Load()
	util.CheckError(err)

	parameter := map[string]string{
		"q":       "Natural Disaster Canada",
		"tbm":     "nws",
		"api_key": os.Getenv("NEWS_API"),
	}

	search := g.NewGoogleSearch(parameter, os.Getenv("NEWS_API"))
	results, err := search.GetJSON()
	util.CheckError(err)
	news_results := results["news_results"].([]interface{})

	return c.JSON(http.StatusOK, news_results)
}
