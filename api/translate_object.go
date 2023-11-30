package api

import (
	"context"
	"os"
	"sync"

	"github.com/labstack/echo/v4"
	"golang.org/x/text/language"

	"encoding/json"

	"cloud.google.com/go/translate"
	"oursos.com/packages/redis"
	"oursos.com/packages/util"
)

func TranslateObject(c echo.Context) error {

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./translate.json")

	ctx := context.Background()

	client, err := translate.NewClient(ctx)
	util.CheckError(err)

	defer client.Close()

	json_map := make(map[string]interface{})

	errEnc := json.NewDecoder(c.Request().Body).Decode(&json_map)
	util.CheckError(errEnc)

	inputObject := json_map["translateObject"].(map[string]interface{})
	lang := json_map["lang"].(string)
	targetLang := language.MustParse(lang).String()

	redis_client := redis.Client()
	redis_ctx := context.Background()

	exists, err := redis_client.Exists(redis_ctx, "translateobject_"+targetLang).Result()
	util.CheckError(err)
	if exists == 1 {
		translate_obj_json := redis_client.Get(redis_ctx, "translateobject_"+targetLang).Val()
		err = json.Unmarshal([]byte(translate_obj_json), &inputObject)
		util.CheckError(err)
		println("redis")
	} else {
		var wg sync.WaitGroup // WaitGroup to synchronize goroutines
		var mutex sync.Mutex
		for key, value := range inputObject {
			if subMap, ok := value.(map[string]interface{}); ok {
				wg.Add(1)
				go func(subMap map[string]interface{}, key string) {
					defer wg.Done() // Decrement the WaitGroup when the goroutine completes
					mutex.Lock()
					newSubMap := make(map[string]interface{})
					for subKey, subValue := range subMap {
						resp, err := client.Translate(ctx, []string{subValue.(string)}, language.Make(targetLang), nil)
						util.CheckError(err)
						translatedText := resp[0].Text
						newSubMap[subKey] = translatedText
					}
					inputObject[key] = newSubMap
					mutex.Unlock()
				}(subMap, key)
			}
		}
		wg.Wait()

		translate_obj_json, err := json.Marshal(inputObject)
		util.CheckError(err)
		rediserr := redis_client.Set(redis_ctx, "translateobject_"+targetLang, translate_obj_json, 0).Err()
		util.CheckError(rediserr)
	}

	return c.JSON(200, inputObject)
}

// func changeValueForKey(data map[string]interface{}, keyToChange, newValue string) {
// 	for key, value := range data {
// 		if subMap, ok := value.(map[string]interface{}); ok {
// 			changeValueForKey(subMap, keyToChange, newValue)
// 		} else if key == keyToChange {
// 			data[key] = newValue
// 		}
// 	}
// }
