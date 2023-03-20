package gpt35turbo

import (
	"context"
	"errors"
	"log"
	"smarthome_ai_bot/clients"
	"smarthome_ai_bot/entities"
	"time"

	"github.com/sashabaranov/go-openai"
)

type clientImpl struct {
	openAIClient *openai.Client
}

type Config struct {
	OpenAIKey string
}

func NewClient(cfg *Config) (clients.PromptInterface, error) {
	if cfg == nil {
		return nil, errors.New("missing parameter: cfg")
	}

	if cfg.OpenAIKey == "" {
		return nil, errors.New("missing parameter: cfg.OpenAIKey")
	}

	aiClient := openai.NewClient(cfg.OpenAIKey)

	return &clientImpl{
		openAIClient: aiClient,
	}, nil
}

func messagesToOpenAIChatCompletionMessages(messages []*entities.Message) []openai.ChatCompletionMessage {
	openAIChatCompletionMessages := make([]openai.ChatCompletionMessage, len(messages))

	for i, message := range messages {
		openAIChatCompletionMessages[i] = openai.ChatCompletionMessage{
			Content: message.Content,
		}

		switch message.Role {
		case entities.RoleSystem:
			openAIChatCompletionMessages[i].Role = openai.ChatMessageRoleSystem
		case entities.RoleUser:
			openAIChatCompletionMessages[i].Role = openai.ChatMessageRoleUser
		case entities.RoleBot:
			openAIChatCompletionMessages[i].Role = openai.ChatMessageRoleAssistant
		}
	}

	return openAIChatCompletionMessages
}

func (client *clientImpl) RequestNextMessage(ctx context.Context, messages []*entities.Message) (*entities.Message, error) {
	reqCtx, cancelFnc := context.WithTimeout(ctx, time.Second*5)
	defer cancelFnc()

	resp, err := client.openAIClient.CreateChatCompletion(reqCtx, openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: messagesToOpenAIChatCompletionMessages(messages),
	})
	if err != nil {
		log.Printf("ChatCompletion error: %v\n", err)

		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, errors.New("no choices returned from OpenAI")
	}

	return &entities.Message{
		Role:    entities.RoleBot,
		Content: resp.Choices[0].Message.Content,
	}, nil
}
