package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"

	"cloud.google.com/go/translate"
	"github.com/labstack/echo/v4"
	"golang.org/x/text/language"
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

	var wg sync.WaitGroup // WaitGroup to synchronize goroutines

	var mutex sync.Mutex // Mutex to protect the languages slice

	for _, lang := range langs {
		wg.Add(1) // Increment the WaitGroup for each goroutine
		go func(lang translate.Language) {
			defer wg.Done() // Decrement the WaitGroup when the goroutine completes

			transName := TranslateText(lang.Name, lang.Tag.String())
			transtag := TranslateText(lang.Tag.String(), lang.Tag.String())

			// Safely append results to the languages slice
			mutex.Lock()
			languages = append(languages, LanguagePreference{Name: transName, Tag: transtag})
			mutex.Unlock()
		}(lang)
	}

	wg.Wait()

	return c.JSON(http.StatusOK, languages)
}

// Return a list of all supported languages in English
func ListLanguagesEnglish(c echo.Context) error {

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./translate.json")
	targetLanguage := "en"
	ctx := context.Background()
	// targetLanguage := "en"
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

	var wg sync.WaitGroup // WaitGroup to synchronize goroutines

	var mutex sync.Mutex // Mutex to protect the languages slice

	for _, lang := range langs {
		wg.Add(1) // Increment the WaitGroup for each goroutine
		go func(lang translate.Language) {
			defer wg.Done() // Decrement the WaitGroup when the goroutine completes
			mutex.Lock()
			languages = append(languages, LanguagePreference{Name: lang.Name, Tag: lang.Tag.String()})
			mutex.Unlock()
		}(lang)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	return c.JSON(http.StatusOK, languages)
}
