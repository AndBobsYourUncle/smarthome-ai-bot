package handlers

import (
	"context"
	"errors"
	"smarthome_ai_bot/clients"
	smarthomeaibotapiv1 "smarthome_ai_bot/gen/go/proto"
)

type grpcV1 struct {
	promptClient clients.PromptInterface
}

type Config struct {
	PromptClient clients.PromptInterface
}

func NewGrpcV1(cfg *Config) (*grpcV1, error) {
	if cfg == nil {
		return nil, errors.New("missing parameter: cfg")
	}

	if cfg.PromptClient == nil {
		return nil, errors.New("missing parameter: cfg.PromptClient")
	}

	return &grpcV1{
		promptClient: cfg.PromptClient,
	}, nil
}

func (s *grpcV1) SendPrompt(
	ctx context.Context,
	req *smarthomeaibotapiv1.SendPromptRequest,
) (*smarthomeaibotapiv1.SendPromptResponse, error) {
	response, err := s.promptClient.SendPrompt(ctx, req.Prompt)
	if err != nil {
		return nil, err
	}

	return &smarthomeaibotapiv1.SendPromptResponse{Response: response}, nil
}
