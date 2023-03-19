package handlers

import (
	"context"
	"errors"
	"smarthome_ai_bot/bot"
	smarthomeaibotapiv1 "smarthome_ai_bot/gen/go/proto"
)

type grpcV1 struct {
	aiBot bot.Interface
}

type Config struct {
	AIBot bot.Interface
}

func NewGrpcV1(cfg *Config) (*grpcV1, error) {
	if cfg == nil {
		return nil, errors.New("missing parameter: cfg")
	}

	if cfg.AIBot == nil {
		return nil, errors.New("missing parameter: cfg.AIBot")
	}

	return &grpcV1{
		aiBot: cfg.AIBot,
	}, nil
}

func (s *grpcV1) SendPrompt(
	ctx context.Context,
	req *smarthomeaibotapiv1.SendPromptRequest,
) (*smarthomeaibotapiv1.SendPromptResponse, error) {
	response, err := s.aiBot.GetResponseToUserMessage(ctx, req.Prompt)
	if err != nil {
		return nil, err
	}

	return &smarthomeaibotapiv1.SendPromptResponse{Response: response}, nil
}
