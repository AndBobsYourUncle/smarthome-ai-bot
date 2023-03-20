package bot

import (
	"context"
	"errors"
	"log"
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
		{
			Role:    entities.RoleSystem,
			Content: "Jarvis has now forgotten everything it knew about the living room temperature sensor.",
		},
	}

	return &botImpl{
		initialPromptContext: initialPromptContext,
		promptClient:         cfg.PromptClient,
	}, nil
}

func (bot *botImpl) GetResponseToUserMessage(ctx context.Context, userMessage string) (string, error) {
	messagesToSend := append(bot.initialPromptContext, &entities.Message{
		Role:    entities.RoleUser,
		Content: userMessage,
	})

	// send the prompt to the prompt client
	log.Printf("Sending user message: %v", messagesToSend[len(messagesToSend)-1].Content)

	response, err := bot.promptClient.RequestNextMessage(ctx, messagesToSend)
	if err != nil {
		return "", err
	}

	log.Printf("Bot response: %v", response.Content)

	// if the response contains a command, then we need to execute it
	if strings.HasPrefix(response.Content, "\"\"\"") && strings.HasSuffix(response.Content, "\"\"\"") {
		systemMessage, sysErr := bot.executeCommand(response.Content)
		if sysErr != nil {
			return "", sysErr
		}

		log.Printf("System response: %v", systemMessage.Content)

		messagesToSend = append(messagesToSend, systemMessage)

		// send the prompt to the prompt client
		log.Printf("Sending system message: %v", messagesToSend[len(messagesToSend)-1].Content)

		response, err = bot.promptClient.RequestNextMessage(ctx, messagesToSend)
		if err != nil {
			return "", err
		}

		log.Printf("Bot response: %v", response.Content)
	}

	return response.Content, nil
}

func (bot *botImpl) executeCommand(command string) (*entities.Message, error) {
	// extract the command from the string
	command = strings.TrimPrefix(command, "\"\"\"")
	command = strings.TrimSuffix(command, "\"\"\"")

	// first word of the command is the command name
	if len(command) == 0 {
		return nil, errors.New("command is empty")
	}

	commandName := strings.Split(command, " ")[0]

	if commandName == "query" {
		entityName := strings.TrimPrefix(command, "query ")

		switch entityName {
		case "living_room_temperature_sensor":
			return &entities.Message{
				Role:    entities.RoleSystem,
				Content: "75 degrees",
			}, nil
		case "front_door_lock":
			return &entities.Message{
				Role:    entities.RoleSystem,
				Content: "locked",
			}, nil
		case "living_room_light":
			return &entities.Message{
				Role:    entities.RoleSystem,
				Content: "on",
			}, nil
		case "living_room_thermostat":
			return &entities.Message{
				Role:    entities.RoleSystem,
				Content: "72 degrees",
			}, nil
		}
	}

	return &entities.Message{
		Role:    entities.RoleSystem,
		Content: "system executed command successfully",
	}, nil
}
