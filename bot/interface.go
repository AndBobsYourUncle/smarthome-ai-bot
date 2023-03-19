package bot

import "context"

type Interface interface {
	SendPrompt(ctx context.Context, prompt string) (string, error)
}
