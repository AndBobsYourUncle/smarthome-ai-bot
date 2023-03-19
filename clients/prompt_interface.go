package clients

import (
	"context"
	"smarthome_ai_bot/entities"
)

type PromptInterface interface {
	RequestNextMessage(ctx context.Context, messages []*entities.Message) (*entities.Message, error)
}
