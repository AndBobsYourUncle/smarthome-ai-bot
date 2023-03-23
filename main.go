package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"smarthome_ai_bot/bot"
	"smarthome_ai_bot/clients/gpt35turbo"
	"smarthome_ai_bot/clients/home_assistant_api"
	"smarthome_ai_bot/handlers"
)

const (
	// defaultPort is the default port that the server will listen on
	defaultPort = "8080"
)

func main() {
	gpt35turboClient, err := gpt35turbo.NewClient(&gpt35turbo.Config{
		OpenAIKey: os.Getenv("OPENAI_API_KEY"),
	})
	if err != nil {
		log.Fatalf("failed to create gpt35turbo client: %v", err)
	}

	homeAssistantClient, err := home_assistant_api.NewClient(&home_assistant_api.Config{
		ApiHost:     os.Getenv("HOME_ASSISTANT_API_HOST"),
		BearerToken: os.Getenv("HOME_ASSISTANT_BEARER_TOKEN"),
	})
	if err != nil {
		log.Fatalf("failed to create home assistant client: %v", err)
	}

	aiBot, err := bot.NewBot(&bot.Config{
		PromptClient:    gpt35turboClient,
		SmarthomeClient: homeAssistantClient,
		UserShortName:   os.Getenv("USER_SHORT_NAME"),
	})
	if err != nil {
		log.Fatalf("failed to create ai bot: %v", err)
	}

	httpHandler, err := handlers.NewHTTP(&handlers.Config{
		AIBot:          aiBot,
		SpeechEndpoint: os.Getenv("SPEECH_ENDPOINT"),
	})
	if err != nil {
		log.Fatalf("failed to create grpc v1 handler: %v", err)
	}

	http.HandleFunc("/get_prompt_response", httpHandler.GetPromptResponse)

	ctx := context.Background()

	aiBot.CleanMemoryOnTimer(ctx)

	defaultPortEnv := os.Getenv("PORT")

	if defaultPortEnv == "" {
		defaultPortEnv = defaultPort
	}

	err = http.ListenAndServe(":"+defaultPortEnv, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	ctx.Done()
}
