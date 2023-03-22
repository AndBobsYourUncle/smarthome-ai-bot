package clients

import "context"

type SmarthomeAPI interface {
	QueryEntity(ctx context.Context, entityID string) (string, error)
	PerformService(ctx context.Context, service, entityID string) (string, error)
}
