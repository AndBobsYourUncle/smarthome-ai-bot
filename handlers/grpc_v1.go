package handlers

import (
	"context"
	"errors"
	smarthomeaibotapiv1 "smarthome_ai_bot/gen/go/proto"
)

type grpcV1 struct {
	// Your server dependencies here...
}

type Config struct {
	// Your server configuration options here...
}

func NewGrpcV1(cfg *Config) (*grpcV1, error) {
	if cfg == nil {
		return nil, errors.New("missing parameter: cfg")
	}

	return &grpcV1{}, nil
}

func (s *grpcV1) SendPrompt(
	ctx context.Context,
	req *smarthomeaibotapiv1.SendPromptRequest,
) (*smarthomeaibotapiv1.SendPromptResponse, error) {
	return &smarthomeaibotapiv1.SendPromptResponse{Response: "Hello, world!"}, nil
}
