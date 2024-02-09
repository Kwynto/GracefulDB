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

func CallSQuery(ctx context.Context, g gs.GracefulServiceClient, text string) (*gs.Response, error) {
	request := &gs.Request{
		Instruction: text,
	}
	r, err := g.SQuery(ctx, request)
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

	var qrys = []string{
		`Errorable Query!`,
	}

	for i1, v1 := range qrys {
		client := gs.NewGracefulServiceClient(conn)
		r, err := CallSQuery(context.Background(), client, v1)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Response Text %d: %s\n", i1, r.Message)
		}
	}
}
