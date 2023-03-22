package bot

import "context"

type Interface interface {
	GetResponseToUserMessage(ctx context.Context, userMessage string) (string, error)
	CleanMemoryOnTimer(ctx context.Context)
}
