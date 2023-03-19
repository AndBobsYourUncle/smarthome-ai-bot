package bot

import (
	"context"
	"errors"
	"os"
	"smarthome_ai_bot/clients"
	"smarthome_ai_bot/entities"
	"strings"
)

type botImpl struct {
	initialPromptContext []*entities.Message
	promptClient         clients.PromptInterface
}

type Config struct {
	UserShortName string
	PromptClient  clients.PromptInterface
}

func NewBot(cfg *Config) (Interface, error) {
	if cfg == nil {
		return nil, errors.New("missing parameter: cfg")
	}

	if cfg.UserShortName == "" {
		return nil, errors.New("missing parameter: cfg.UserShortName")
	}

	if cfg.PromptClient == nil {
		return nil, errors.New("missing parameter: cfg.PromptClient")
	}

	// read initial context from file
	initialContext, err := os.ReadFile("initial_context.txt")
	if err != nil {
		return nil, err
	}

	// replace all occurrences of "[USERS_SHORT_NAME]" with the user's short name
	initialContext = []byte(strings.ReplaceAll(string(initialContext), "[USERS_SHORT_NAME]", cfg.UserShortName))

	initialPromptContext := []*entities.Message{
		{
			Role:    entities.RoleSystem,
			Content: string(initialContext),
		},
		{
			Role:    entities.RoleUser,
			Content: "What is the temperature in the living room?",
		},
		{
			Role:    entities.RoleBot,
			Content: "\"\"\"query living_room_temperature_sensor\"\"\"",
		},
		{
			Role:    entities.RoleSystem,
			Content: "72 degrees",
		},
		{
			Role:    entities.RoleBot,
			Content: "The living room temperature appears to be 72 degrees.",
		},
	}

	return &botImpl{
		initialPromptContext: initialPromptContext,
		promptClient:         cfg.PromptClient,
	}, nil
}

func (client *botImpl) SendPrompt(ctx context.Context, prompt string) (string, error) {
	messagesToSend := append(client.initialPromptContext, &entities.Message{
		Role:    entities.RoleUser,
		Content: prompt,
	})

	// send the prompt to the prompt client
	response, err := client.promptClient.RequestNextMessage(ctx, messagesToSend)
	if err != nil {
		return "", err
	}

	return response.Content, nil
}
