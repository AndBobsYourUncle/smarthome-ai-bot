package handlers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"smarthome_ai_bot/bot"
)

type httpImpl struct {
	aiBot bot.Interface
}

type Config struct {
	AIBot bot.Interface
}

func NewHTTP(cfg *Config) (*httpImpl, error) {
	if cfg == nil {
		return nil, errors.New("missing parameter: cfg")
	}

	if cfg.AIBot == nil {
		return nil, errors.New("missing parameter: cfg.AIBot")
	}

	return &httpImpl{
		aiBot: cfg.AIBot,
	}, nil
}

func (s *httpImpl) GetPromptResponse(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	v := r.URL.Query()

	prompt := v["prompt"][0]

	response, err := s.aiBot.GetResponseToUserMessage(ctx, prompt)
	if err != nil {
		w.WriteHeader(500)
		_, _ = io.WriteString(w, "error: "+err.Error())

		return
	}

	fmt.Fprintf(w, "%s", response)
}
