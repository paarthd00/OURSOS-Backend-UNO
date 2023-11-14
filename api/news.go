package api

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	g "github.com/serpapi/google-search-results-golang"
	"oursos.com/packages/redis"
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
	client := redis.Client()
	ctx := context.Background()

	exists, err := client.Exists(ctx, "news").Result()
	util.CheckError(err)

	var redis_news []interface{}

	if exists == 1 {
		redis_news_json := client.Get(ctx, "news").Val()
		err = json.Unmarshal([]byte(redis_news_json), &redis_news)
		util.CheckError(err)
		println("redis")
	} else {
		parameter := map[string]string{
			"q":       "Natural Disaster Canada",
			"tbm":     "nws",
			"api_key": os.Getenv("NEWS_API"),
		}
		search := g.NewGoogleSearch(parameter, os.Getenv("NEWS_API"))
		results, err := search.GetJSON()
		util.CheckError(err)
		news_results := results["news_results"].([]interface{})

		news_results_json, err := json.Marshal(news_results)
		util.CheckError(err)
		rediserr := client.Set(ctx, "news", news_results_json, 0).Err()
		util.CheckError(rediserr)
	}

	return c.JSON(http.StatusOK, redis_news)
}
