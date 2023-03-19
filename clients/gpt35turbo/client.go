package gpt35turbo

import (
	"context"
	"errors"
	"log"
	"os"
	"smarthome_ai_bot/clients"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type clientImpl struct {
	openAIKey     string
	userShortName string
}

type Config struct {
	OpenAIKey     string
	UserShortName string
}

func NewClient(cfg *Config) (clients.PromptInterface, error) {
	if cfg == nil {
		return nil, errors.New("missing parameter: cfg")
	}

	if cfg.OpenAIKey == "" {
		return nil, errors.New("missing parameter: cfg.OpenAIKey")
	}

	if cfg.UserShortName == "" {
		return nil, errors.New("missing parameter: cfg.UserShortName")
	}

	return &clientImpl{
		openAIKey:     cfg.OpenAIKey,
		userShortName: cfg.UserShortName,
	}, nil
}

func (client *clientImpl) SendPrompt(ctx context.Context, prompt string) (string, error) {
	// read initial context from file
	initialContext, err := os.ReadFile("initial_context.txt")
	if err != nil {
		return "", err
	}

	// replace all occurrences of "[USERS_SHORT_NAME]" with the user's short name
	initialContext = []byte(strings.ReplaceAll(string(initialContext), "[USERS_SHORT_NAME]", client.userShortName))

	aiClient := openai.NewClient(client.openAIKey)
	resp, err := aiClient.CreateChatCompletion(
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
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		log.Printf("ChatCompletion error: %v\n", err)

		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("no choices returned from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}
