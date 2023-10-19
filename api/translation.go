package api

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/labstack/echo/v4"

	"encoding/json"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"oursos.com/packages/util"
)

func TranslateText(text string, trg string) string {
	ctx := context.Background()

	client, err := translate.NewClient(ctx)
	util.CheckError(err)

	defer client.Close()

	inputText := text
	lang := trg
	targetLang := language.MustParse(lang)

	target := targetLang.String()

	resp, err := client.Translate(ctx, []string{inputText}, language.Make(target), nil)

	util.CheckError(err)

	translatedText := resp[0].Text

	return translatedText
}

func Translate(c echo.Context) error {

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./translate.json")

	ctx := context.Background()

	client, err := translate.NewClient(ctx)
	util.CheckError(err)

	defer client.Close()

	json_map := make(map[string]interface{})

	fmt.Println(reflect.TypeOf(c.Request().Body))
	errEnc := json.NewDecoder(c.Request().Body).Decode(&json_map)
	util.CheckError(errEnc)

	inputText := json_map["text"].(string)
	lang := json_map["lang"].(string)
	targetLang := language.MustParse(lang)

	target := targetLang.String()

	resp, err := client.Translate(ctx, []string{inputText}, language.Make(target), nil)

	util.CheckError(err)

	translatedText := resp[0].Text

	return c.JSON(200, translatedText)
}
