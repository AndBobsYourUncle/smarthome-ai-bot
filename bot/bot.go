package bot

import (
	"context"
	"errors"
	"log"
	"os"
	"regexp"
	"smarthome_ai_bot/clients"
	"smarthome_ai_bot/entities"
	"strings"
)

type botImpl struct {
	initialPromptContext []*entities.Message
	messageHistory       []*entities.Message
	promptClient         clients.PromptInterface
	smarthomeClient      clients.SmarthomeAPI
}

type Config struct {
	UserShortName   string
	PromptClient    clients.PromptInterface
	SmarthomeClient clients.SmarthomeAPI
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

	if cfg.SmarthomeClient == nil {
		return nil, errors.New("missing parameter: cfg.SmarthomeClient")
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
			Content: "\"\"\"query sensor.living_room_sensor_air_temperature\"\"\"",
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
			Content: "Jarvis never remembers device sensor values, and must query the system every time.",
		},
	}

	return &botImpl{
		initialPromptContext: initialPromptContext,
		promptClient:         cfg.PromptClient,
		smarthomeClient:      cfg.SmarthomeClient,
	}, nil
}

func (bot *botImpl) getMessagesToSend() []*entities.Message {
	// return a new slice of messages that contains the initial prompt context and the message history
	messagesToSend := make([]*entities.Message, len(bot.initialPromptContext)+len(bot.messageHistory))
	copy(messagesToSend, bot.initialPromptContext)
	copy(messagesToSend[len(bot.initialPromptContext):], bot.messageHistory)

	return messagesToSend
}

const stringCommandRegex = `\"\"\".*\"\"\"`

func extractStringCommand(s string) string {
	r := regexp.MustCompile(stringCommandRegex)

	// extract the command from the string
	command := r.FindString(s)

	if command == "" {
		return ""
	}

	return command
}

func (bot *botImpl) GetResponseToUserMessage(ctx context.Context, userMessage string) (string, error) {
	messagesToSend := bot.getMessagesToSend()

	lengthOfInitialMessages := len(messagesToSend)

	messagesToSend = append(messagesToSend, &entities.Message{
		Role:    entities.RoleUser,
		Content: userMessage,
	})

	// send the prompt to the prompt client
	log.Printf("Sending user message: %v", messagesToSend[len(messagesToSend)-1].Content)

	response, err := bot.promptClient.RequestNextMessage(ctx, messagesToSend)
	if err != nil {
		return "", err
	}

	botResponse := response.Content

	stringCommand := extractStringCommand(response.Content)

	if stringCommand != "" {
		botResponse = stringCommand
	}

	if len(botResponse) < len(response.Content) {
		log.Printf("Bot's original response: %v", response.Content)
	}

	log.Printf("Bot response: %v", botResponse)

	messagesToSend = append(messagesToSend, &entities.Message{
		Role:    entities.RoleBot,
		Content: botResponse,
	})

	// if the response contains a command, then we need to execute it
	if stringCommand != "" {
		systemMessage, sysErr := bot.executeCommand(ctx, stringCommand)
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

		messagesToSend = append(messagesToSend, &entities.Message{
			Role:    entities.RoleBot,
			Content: response.Content,
		})

		log.Printf("Bot response: %v", response.Content)
	}

	// add all new messages to the message history if we are successfully able to get a final response
	bot.messageHistory = append(bot.messageHistory, messagesToSend[lengthOfInitialMessages:]...)

	return response.Content, nil
}

func (bot *botImpl) executeCommand(ctx context.Context, command string) (*entities.Message, error) {
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

		entityResponse, err := bot.smarthomeClient.QueryEntity(ctx, entityName)
		if err != nil {
			return &entities.Message{
				Role:    entities.RoleSystem,
				Content: "system failed to execute command",
			}, nil
		}

		return &entities.Message{
			Role:    entities.RoleSystem,
			Content: entityResponse,
		}, nil
	} else {
		entityName := strings.TrimPrefix(command, commandName)

		// remove the first space
		entityName = strings.TrimPrefix(entityName, " ")

		// extract a value if there is a space and then a value
		value := ""

		if strings.Contains(entityName, " ") {
			value = strings.Split(entityName, " ")[1]
			entityName = strings.Split(entityName, " ")[0]
		}

		response, err := bot.smarthomeClient.PerformService(ctx, commandName, entityName, value)
		if err != nil {
			return &entities.Message{
				Role:    entities.RoleSystem,
				Content: "system failed to execute command",
			}, nil
		}

		return &entities.Message{
			Role:    entities.RoleSystem,
			Content: response,
		}, nil
	}
}
