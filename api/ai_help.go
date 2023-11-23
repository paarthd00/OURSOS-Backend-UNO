package api

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"

// 	"github.com/labstack/echo/v4"
// 	"oursos.com/packages/util"
// )

// func AskOpenAI(c echo.Context) error {
// 	url := "https://api.openai.com/v1/engines/davinci-codex/completions"
// 	var jsonStr = []byte(`{"prompt":"` + "I have a fire near me" + `", "max_tokens": 60}`)
// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

// 	util.CheckError(err)
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", "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")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	var bodyJSON interface{}
// 	json.NewDecoder(resp.Body).Decode(&bodyJSON)
// 	fmt.Println("response Body:", bodyJSON)

// 	return c.JSON(http.StatusOK, bodyJSON)
// }
