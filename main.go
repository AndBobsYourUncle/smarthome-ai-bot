package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

func main() {
	// get token from env var
	token := os.Getenv("OPENAI_API_KEY")

	// get user's short name from env var
	userShortName := os.Getenv("USER_SHORT_NAME")

	// read initial context from file
	initialContext, err := os.ReadFile("initial_context.txt")
	if err != nil {
		panic(err)
	}

	// replace all occurrences of "[USERS_SHORT_NAME]" with the user's short name
	initialContext = []byte(strings.ReplaceAll(string(initialContext), "[USERS_SHORT_NAME]", userShortName))

	client := openai.NewClient(token)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: string(initialContext),
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "What is the temperature in the living room?",
				},
				{
					Role:    openai.ChatMessageRoleAssistant,
					Content: "\"\"\"query living_room_temperature_sensor\"\"\"",
				},
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "72 degrees",
				},
				{
					Role:    openai.ChatMessageRoleAssistant,
					Content: "The living room temperature appears to be 72 degrees.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Who created you?",
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
}
