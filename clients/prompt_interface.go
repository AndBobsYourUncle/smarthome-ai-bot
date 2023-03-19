package clients

import "context"

type PromptInterface interface {
	SendPrompt(ctx context.Context, prompt string) (string, error)
}
