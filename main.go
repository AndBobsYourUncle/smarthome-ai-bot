package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"smarthome_ai_bot/clients/gpt35turbo"
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

	gpt35turboClient, err := gpt35turbo.NewClient(&gpt35turbo.Config{
		OpenAIKey:     os.Getenv("OPENAI_API_KEY"),
		UserShortName: os.Getenv("USER_SHORT_NAME"),
	})
	if err != nil {
		log.Fatalf("failed to create gpt35turbo client: %v", err)
	}

	handler, err := handlers.NewGrpcV1(&handlers.Config{
		PromptClient: gpt35turboClient,
	})

	smarthomeaibotapiv1.RegisterSmarthomeAIBotAPIServer(grpcServer, handler)

	ctx := context.Background()

	err = runServer(ctx, grpcServer, defaultPort)
	if err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
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
