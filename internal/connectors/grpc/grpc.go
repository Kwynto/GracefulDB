package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	gs "github.com/Kwynto/GracefulDB/internal/connectors/grpc/proto/graceful_service"
	"google.golang.org/grpc"

	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
)

type tMessageServer struct {
	gs.UnimplementedGracefulServiceServer
}

var address string

var messageServer tMessageServer
var server *grpc.Server

func (tMessageServer) Query(ctx context.Context, r *gs.Request) (*gs.Response, error) {
	slog.Debug("Request received", slog.String("request", r.Message))

	// TODO: There should be request processing here
	response := &gs.Response{
		Message: "There should be a response from the processed request.",
	}
	slog.Debug("Response sent", slog.String("response", response.Message))

	return response, nil
}

func Start(cfg *config.Config) {
	address = fmt.Sprintf("%s:%s", cfg.GrpcConnector.Address, cfg.GrpcConnector.Port)
	listen, err := net.Listen("tcp", address)
	if err != nil {
		slog.Error("Failed to start listener", slog.String("err", err.Error()))
		return
	}

	server = grpc.NewServer()
	gs.RegisterGracefulServiceServer(server, messageServer)

	server.Serve(listen)
}

func Shutdown(ctx context.Context, c *closer.Closer) {
	server.Stop()
	c.Done()
}
