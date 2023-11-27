package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/sashabaranov/go-openai"
	"oursos.com/packages/util"
)

type ChatRequest struct {
	Message string `json:"message"`
}

func ChatHandler(c echo.Context) error {
	err := godotenv.Load()
	util.CheckError(err)

	var userInput ChatRequest

	user_input_err := json.NewDecoder(c.Request().Body).Decode(&userInput)
	if user_input_err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	client := openai.NewClient(os.Getenv("OPEN_AI"))
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userInput.Message + " ###DEV please give concise instructions for emergency in less than 100 characters more than 50 characters",
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"response": resp.Choices[0].Message.Content})
}
