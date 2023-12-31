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

	var qrys = []string{
		`Errorable Query!`,
		`{}`,
		`{"action":""}`,
		`{"action":"auth", "secret":{}}`,
		`{"action":"auth", "secret":{"login":"root", "password":"toor"}}`,
		`{"action":"auth", "secret":{"login":"root", "password":"toor", "queryid":"any-id"}}`,
		`{"action":"read", "secret":{}, "db":""}`,
		`{"action":"store", "secret":{}, "db":"", "table":""}`,
		`{"action":"delete", "secret":{}, "db":"", "table":"", "fields":{}}`,
		`{"action":"manage", "secret":{}, "db":"", "table":"", "fields":{}, "data":[]}`,
		`{"action":"auth", "secret":{}, "db":"", "table":"", "fields":{}, "data":["Errorable Query!"]}`,
		`{"action":"read", "secret":{}, "db":"", "table":"", "fields":{}, "data":[{}]}`,
		`{"action":"store", "secret":{}, "db":"", "table":"", "fields":{}, "data":[{},{}]}`,
		`{"action":"delete", "secret":{}, "db":"", "table":"", "fields":{}, "data":[{"name":"Name"},{"name":"Name","city":"Moscow"}]}`,
		`{"action":"manage", "secret":{}, "db":"", "table":"", "fields":{}, "data":[{"name":"Name"},{"name":"Name","city":"Moscow","sub":""}]}`,
		`{"action":"auth", "secret":{}, "db":"", "table":"", "fields":{}, "data":[{"name":"Name"},{"name":"Name","city":"Moscow","sub":"","age":20}]}`,
		`{"action":"read", "secret":{}, "db":"", "table":"", "fields":{"name":"Name","city":"Moscow","sub":"","age":20}, "data":[{"name":"Name"},{"name":"Name","city":"Moscow","sub":"","age":20}]}`,
		`{"action":"store", "secret":{}, "db":"", "table":"", "fields":{"name":"Name","city":"Moscow","sub":"","age":20}, "data":[{"name":"Name"},{"name":"Name","city":"Moscow","sub":"","age":20}]}`,
		`{"action":"delete", "secret":{}, "db":"", "table":"", "fields":{"name":"Name","city":"Moscow","sub":"","age":20}, "data":[{"name":"Name"},{"name":"Name","city":"Moscow","sub":"","age":20}]}`,
	}

	for i1, v1 := range qrys {
		client := gs.NewGracefulServiceClient(conn)
		r, err := CallVQuery(context.Background(), client, v1)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Response Text %d: %s\n", i1, r.Message)
		}
	}
}
