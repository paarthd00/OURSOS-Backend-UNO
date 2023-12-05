package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	"cloud.google.com/go/translate"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	g "github.com/serpapi/google-search-results-golang"
	"golang.org/x/text/language"
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

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./translate.json")
	err := godotenv.Load()
	util.CheckError(err)

	lang := c.Param("lang")

	parameter := map[string]string{
		"q":       "Natural Disaster Canada",
		"tbm":     "nws",
		"api_key": os.Getenv("NEWS_API"),
	}
	search := g.NewGoogleSearch(parameter, os.Getenv("NEWS_API"))
	results, err := search.GetJSON()
	util.CheckError(err)
	news_results := results["news_results"].([]interface{})

	redis_client := redis.Client()
	redis_ctx := context.Background()

	exists, redis_err := redis_client.Exists(redis_ctx, "news_"+lang).Result()
	util.CheckError(redis_err)

	ctx := context.Background()

	client, err := translate.NewClient(ctx)
	util.CheckError(err)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	if exists == 1 {
		news_results_json, news_err := redis_client.Get(redis_ctx, "news_"+lang).Result()
		util.CheckError(news_err)
		err = json.Unmarshal([]byte(news_results_json), &news_results)
		util.CheckError(err)
		println("redis")
	} else {
		if lang != "en" {
			var wg sync.WaitGroup
			for _, article := range news_results {
				wg.Add(1)
				go func(article interface{}) {
					defer wg.Done()
					articleMap := article.(map[string]interface{})
					title := articleMap["title"].(string)
					snippet := articleMap["snippet"].(string)

					// Translate title
					titleTranslation, err := client.Translate(ctx, []string{title}, language.Make(lang), nil)
					if err != nil {
						log.Fatalf("Failed to translate text: %v", err)
					}
					articleMap["title"] = titleTranslation[0].Text

					// Translate snippet
					snippetTranslation, err := client.Translate(ctx, []string{snippet}, language.Make(lang), nil)
					if err != nil {
						log.Fatalf("Failed to translate text: %v", err)
					}
					articleMap["snippet"] = snippetTranslation[0].Text
				}(article)
			}
			wg.Wait()
		}
		news_results_json, err := json.Marshal(news_results)
		util.CheckError(err)
		redis_err := redis_client.Set(redis_ctx, "news_"+lang, news_results_json, 0).Err()
		util.CheckError(redis_err)
	}

	return c.JSON(http.StatusOK, news_results)
}
