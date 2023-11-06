package main

// TODO: Delete this code along with the root folder after debugging.

import (
	"context"
	"fmt"
	"os"

	"github.com/Kwynto/GracefulDB/internal/config"

	gs "github.com/Kwynto/GracefulDB/internal/connectors/grpc/proto/graceful_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CallVQuery(ctx context.Context, g gs.GracefulServiceClient, text string) (*gs.VResponse, error) {
	request := &gs.VRequest{
		Instruction: text,
	}
	r, err := g.VQuery(ctx, request)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func main() {
	// Init config
	configPath := os.Getenv("CONFIG_PATH")
	cfg := config.MustLoad(configPath)
	address := fmt.Sprintf("%s:%s", cfg.GrpcConnector.Address, cfg.GrpcConnector.Port)

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Dial:", err)
		return
	}

	client := gs.NewGracefulServiceClient(conn)
	r, err := CallVQuery(context.Background(), client, "Database Query!")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Response Text:", r.Message)
	}

}
