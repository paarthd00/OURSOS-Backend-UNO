package api

import (
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"encoding/json"

	"oursos.com/packages/util"
)

func Speech(c echo.Context) error {
	enverr := godotenv.Load()
	util.CheckError(enverr)

	url := "https://joj-text-to-speech.p.rapidapi.com/"

	payload := strings.NewReader("{\n    \"input\": {\n        \"text\": \"Mary and Samantha arrived at the bus station early but waited until noon for the bus.\"\n    },\n    \"voice\": {\n        \"languageCode\": \"en-US\",\n        \"name\": \"en-US-News-L\",\n        \"ssmlGender\": \"FEMALE\"\n    },\n    \"audioConfig\": {\n        \"audioEncoding\": \"MP3\"\n    }\n}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-RapidAPI-Key", os.Getenv("AI_SPEACH_KEY"))
	req.Header.Add("X-RapidAPI-Host", "joj-text-to-speech.p.rapidapi.com")

	res, err := http.DefaultClient.Do(req)
	util.CheckError(err)
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return c.JSON(res.StatusCode, map[string]string{"error": "Failed to retrieve speech"})
	}
	var bodyJSON interface{}
	// Read the response body
	json.NewDecoder(res.Body).Decode(&bodyJSON)

	util.CheckError(err)
	// Send the JSON response to the client
	return c.JSON(http.StatusOK, bodyJSON)
}
