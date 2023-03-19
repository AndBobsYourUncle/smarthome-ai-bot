package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	smarthomeaibotapiv1 "smarthome_ai_bot/gen/go/proto"
	"smarthome_ai_bot/handlers"
	"syscall"

	"google.golang.org/grpc"
)

const (
	// defaultPort is the default port that the server will listen on
	defaultPort = 5001
)

func main() {
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)

	handler, err := handlers.NewGrpcV1(&handlers.Config{})

	smarthomeaibotapiv1.RegisterSmarthomeAIBotAPIServer(grpcServer, handler)

	ctx := context.Background()

	err = runServer(ctx, grpcServer, defaultPort)
	if err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

	//// get token from env var
	//token := os.Getenv("OPENAI_API_KEY")
	//
	//// get user's short name from env var
	//userShortName := os.Getenv("USER_SHORT_NAME")
	//
	//// read initial context from file
	//initialContext, err := os.ReadFile("initial_context.txt")
	//if err != nil {
	//	panic(err)
	//}
	//
	//// replace all occurrences of "[USERS_SHORT_NAME]" with the user's short name
	//initialContext = []byte(strings.ReplaceAll(string(initialContext), "[USERS_SHORT_NAME]", userShortName))
	//
	//client := openai.NewClient(token)
	//resp, err := client.CreateChatCompletion(
	//	context.Background(),
	//	openai.ChatCompletionRequest{
	//		Model: openai.GPT3Dot5Turbo,
	//		Messages: []openai.ChatCompletionMessage{
	//			{
	//				Role:    openai.ChatMessageRoleSystem,
	//				Content: string(initialContext),
	//			},
	//			{
	//				Role:    openai.ChatMessageRoleUser,
	//				Content: "What is the temperature in the living room?",
	//			},
	//			{
	//				Role:    openai.ChatMessageRoleAssistant,
	//				Content: "\"\"\"query living_room_temperature_sensor\"\"\"",
	//			},
	//			{
	//				Role:    openai.ChatMessageRoleSystem,
	//				Content: "72 degrees",
	//			},
	//			{
	//				Role:    openai.ChatMessageRoleAssistant,
	//				Content: "The living room temperature appears to be 72 degrees.",
	//			},
	//			{
	//				Role:    openai.ChatMessageRoleUser,
	//				Content: "Who created you?",
	//			},
	//		},
	//	},
	//)
	//
	//if err != nil {
	//	fmt.Printf("ChatCompletion error: %v\n", err)
	//	return
	//}
	//
	//fmt.Println(resp.Choices[0].Message.Content)
}

func runServer(ctx context.Context, srv *grpc.Server, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	runCtx, cancel := context.WithCancel(ctx)

	go watchSignals(runCtx, cancel)

	go func() {
		if err = srv.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
		}
	}()

	log.Printf("!!! grpc server listening on port %d\n", port)

	<-runCtx.Done()

	log.Printf("!!! shutting down grpc server")

	srv.GracefulStop()

	return nil
}

func watchSignals(ctx context.Context, fn context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	sig := <-ch
	log.Printf("!!! received signal, shutting down: %v", sig)
	fn()
}
