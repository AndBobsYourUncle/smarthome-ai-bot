package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"smarthome_ai_bot/bot"
)

type httpImpl struct {
	aiBot          bot.Interface
	speechEndpoint string
}

type Config struct {
	AIBot          bot.Interface
	SpeechEndpoint string
}

func NewHTTP(cfg *Config) (*httpImpl, error) {
	if cfg == nil {
		return nil, errors.New("missing parameter: cfg")
	}

	if cfg.AIBot == nil {
		return nil, errors.New("missing parameter: cfg.AIBot")
	}

	if cfg.SpeechEndpoint == "" {
		return nil, errors.New("missing parameter: cfg.SpeechEndpoint")
	}

	return &httpImpl{
		aiBot:          cfg.AIBot,
		speechEndpoint: cfg.SpeechEndpoint,
	}, nil
}

func (s *httpImpl) GetPromptResponse(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	v := r.URL.Query()

	if _, ok := v["prompt"]; !ok {
		w.WriteHeader(422)
		_, _ = io.WriteString(w, "error: missing prompt query parameter")

		return
	} else if len(v["prompt"]) == 0 {
		w.WriteHeader(422)
		_, _ = io.WriteString(w, "error: prompt query parameter is empty")

		return
	}

	prompt := v["prompt"][0]

	response, err := s.aiBot.GetResponseToUserMessage(ctx, prompt)
	if err != nil {
		w.WriteHeader(500)
		_, _ = io.WriteString(w, "error: "+err.Error())

		return
	}

	postBody, _ := json.Marshal(map[string]string{
		"response": response,
	})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(s.speechEndpoint, "application/json", responseBody)
	if err != nil {
		log.Printf("An Error Occured %v\n", err)
	}
	defer resp.Body.Close()

	fmt.Fprintf(w, "%s", response)
}
