package api

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"sync"

	"github.com/labstack/echo/v4"
	"golang.org/x/text/language"

	"encoding/json"

	"cloud.google.com/go/translate"
	"oursos.com/packages/util"
)

func TranslateObject(c echo.Context) error {

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./translate.json")

	ctx := context.Background()

	client, err := translate.NewClient(ctx)
	util.CheckError(err)

	defer client.Close()

	json_map := make(map[string]interface{})

	fmt.Println(reflect.TypeOf(c.Request().Body))
	errEnc := json.NewDecoder(c.Request().Body).Decode(&json_map)
	util.CheckError(errEnc)

	inputObject := json_map["translateObject"].(map[string]interface{})
	lang := json_map["lang"].(string)
	targetLang := language.MustParse(lang).String()
	var wg sync.WaitGroup // WaitGroup to synchronize goroutines

	var mutex sync.Mutex
	for _, value := range inputObject {

		if subMap, ok := value.(map[string](interface{})); ok {
			for subKey, subValue := range subMap {
				wg.Add(1)
				go func(subKey string, subValue interface{}) {
					defer wg.Done() // Decrement the WaitGroup when the goroutine completes
					mutex.Lock()
					resp, err := client.Translate(ctx, []string{subValue.(string)}, language.Make(targetLang), nil)
					util.CheckError(err)

					translatedText := resp[0].Text
					changeValueForKey(inputObject, subKey, translatedText)
					mutex.Unlock()
				}(subKey, subValue)
			}
		}

	}
	wg.Wait()
	return c.JSON(200, inputObject)
}

func changeValueForKey(data map[string]interface{}, keyToChange, newValue string) {
	for key, value := range data {
		if subMap, ok := value.(map[string]interface{}); ok {
			// If the value is a map (object), recursively change the value
			changeValueForKey(subMap, keyToChange, newValue)
		} else if key == keyToChange {
			// If the key matches the key to change, update the value
			data[key] = newValue
		}
	}
}
