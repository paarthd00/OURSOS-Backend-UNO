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

var staticJSON = map[string]interface{}{
	"dashboard": map[string]interface{}{
		"news":      "News",
		"map":       "Map",
		"getAIHelp": "Get AI Help",
	},
	"map": map[string]interface{}{
		"standard":    "Standard",
		"report":      "Report",
		"shrink":      "Shrink",
		"list":        "List",
		"all":         "All",
		"hazards":     "Hazards",
		"fires":       "Fires",
		"police":      "Police",
		"earthquakes": "EarthQuakes",
		"tsunamis":    "Tsunamis",
		"wildfires":   "Wild Fires",
		"show":        "Show",
		"hybrid":      "Hybrid",
		"default":     "Default",
	},
	"menu": map[string]interface{}{
		"home":     "Home",
		"settings": "Settings",
	},
	"modal": map[string]interface{}{
		"page":          "Page",
		"whatdidyousee": "What did you see",
		"severity":      "Severity",
		"tellusmore":    "Tell Us More",
		"whathappened":  "What Happened",
		"next":          "Next",
		"description":   "Description",
	},
	"settings": map[string]interface{}{
		"changelanguage": "Change Language",
		"updateprofile":  "Update Profile",
		"addfriend":      "Add Friend",
	},
}

func TranslateObject(c echo.Context) error {

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./translate.json")

	ctx := context.Background()

	client, err := translate.NewClient(ctx)
	util.CheckError(err)

	defer client.Close()

	lang := c.Param("lang")

	redis_client := redis.Client()
	redis_ctx := context.Background()

	exists, err := redis_client.Exists(redis_ctx, "translateobject_"+lang).Result()
	util.CheckError(err)
	if exists == 1 {
		translate_obj_json := redis_client.Get(redis_ctx, "translateobject_"+lang).Val()
		err = json.Unmarshal([]byte(translate_obj_json), &staticJSON)
		util.CheckError(err)
		println("redis")
	} else {
		var wg sync.WaitGroup // WaitGroup to synchronize goroutines
		var mutex sync.Mutex
		for key, value := range staticJSON {
			if subMap, ok := value.(map[string]interface{}); ok {
				wg.Add(1)
				go func(subMap map[string]interface{}, key string) {
					defer wg.Done() // Decrement the WaitGroup when the goroutine completes
					mutex.Lock()
					newSubMap := make(map[string]interface{})
					for subKey, subValue := range subMap {
						resp, err := client.Translate(ctx, []string{subValue.(string)}, language.Make(lang), nil)
						util.CheckError(err)
						translatedText := resp[0].Text
						newSubMap[subKey] = translatedText
					}
					staticJSON[key] = newSubMap
					mutex.Unlock()
				}(subMap, key)
			}
		}
		wg.Wait()

		translate_obj_json, err := json.Marshal(staticJSON)
		util.CheckError(err)
		rediserr := redis_client.Set(redis_ctx, "translateobject_"+lang, translate_obj_json, 0).Err()
		util.CheckError(rediserr)
	}

	return c.JSON(200, staticJSON)
}
