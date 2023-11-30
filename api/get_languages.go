package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"cloud.google.com/go/translate"
	"github.com/labstack/echo/v4"
	"golang.org/x/text/language"
	"oursos.com/packages/redis"
	"oursos.com/packages/util"
)

type LanguagePreference struct {
	Name string `json:"name"`
	Tag  string `json:"tag"`
}

// Gets a list of all supported languages in original Language
func ListSupportedLanguages(c echo.Context) error {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./translate.json")
	targetLanguage := "en"
	ctx := context.Background()

	lang, err := language.Parse(targetLanguage)
	if err != nil {
		return fmt.Errorf("language.Parse: %w", err)
	}

	client, err := translate.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("translate.NewClient: %w", err)
	}
	defer client.Close()

	langs, err := client.SupportedLanguages(ctx, lang)
	if err != nil {
		return fmt.Errorf("SupportedLanguages: %w", err)
	}
	var languages []LanguagePreference

	redis_client := redis.Client()
	redis_ctx := context.Background()

	exists, err := redis_client.Exists(redis_ctx, "languages").Result()
	util.CheckError(err)
	if exists == 1 {
		languages_json := redis_client.Get(redis_ctx, "languages").Val()
		err = json.Unmarshal([]byte(languages_json), &languages)
		util.CheckError(err)
		println("redis")

	} else {
		var wg sync.WaitGroup // WaitGroup to synchronize goroutines

		var mutex sync.Mutex // Mutex to protect the languages slice

		for _, lang := range langs {
			wg.Add(1) // Increment the WaitGroup for each goroutine
			go func(lang translate.Language) {
				defer wg.Done() // Decrement the WaitGroup when the goroutine completes

				transName := TranslateText(lang.Name, lang.Tag.String())
				// transtag := TranslateText(lang.Tag.String(), lang.Tag.String())

				// Safely append results to the languages slice
				mutex.Lock()
				languages = append(languages, LanguagePreference{Name: transName, Tag: lang.Tag.String()})
				mutex.Unlock()
			}(lang)
		}
		wg.Wait()
		languages_json, err := json.Marshal(languages)
		util.CheckError(err)
		rediserr := redis_client.Set(ctx, "languages", languages_json, 0).Err()
		util.CheckError(rediserr)
	}
	return c.JSON(http.StatusOK, languages)
}
