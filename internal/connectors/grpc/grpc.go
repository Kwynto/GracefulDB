package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/Kwynto/GracefulDB/internal/analyzers/sqlanalyzer"
	"github.com/Kwynto/GracefulDB/internal/analyzers/vqlanalyzer"
	gs "github.com/Kwynto/GracefulDB/internal/connectors/grpc/proto/graceful_service"
	"google.golang.org/grpc"

	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
	"github.com/Kwynto/GracefulDB/pkg/lib/e"
)

type tMessageServer struct {
	gs.UnimplementedGracefulServiceServer
}

var address string

var messageServer tMessageServer
var grpcServer *grpc.Server

func (tMessageServer) SQuery(ctx context.Context, r *gs.SRequest) (response *gs.SResponse, err error) {
	op := "internal -> connectors -> gRPC -> SQuery"
	defer func() { e.Wrapper(op, err) }()

	slog.Debug("Request received", slog.String("instruction", r.Instruction), slog.String("placeholder", fmt.Sprint(r.Placeholder)))

	// instructionB := []byte(r.Instruction)
	// placeholderB := []byte(r.Placeholder)

	response = &gs.SResponse{
		Message: *sqlanalyzer.Request(&r.Ticket, &r.Instruction, &r.Placeholder),
	}
	slog.Debug("Response sent", slog.String("response", response.Message))

	return response, nil
}

func (tMessageServer) VQuery(ctx context.Context, r *gs.VRequest) (response *gs.VResponse, err error) {
	op := "internal -> connectors -> gRPC -> VQuery"
	defer func() { e.Wrapper(op, err) }()

	slog.Debug("Request received", slog.String("request", r.Instruction))

	instructionB := []byte(r.Instruction)

	response = &gs.VResponse{
		Message: string(*vqlanalyzer.Request(&instructionB)),
	}
	slog.Debug("Response sent", slog.String("response", response.Message))

	return response, nil
}

func Start(cfg *config.Config) {
	address = fmt.Sprintf("%s:%s", cfg.GrpcConnector.Address, cfg.GrpcConnector.Port)
	listen, err := net.Listen("tcp", address)
	if err != nil {
		slog.Error("Failed to start gRPC-listener", slog.String("err", err.Error()))
		return
	}

	grpcServer = grpc.NewServer()
	gs.RegisterGracefulServiceServer(grpcServer, messageServer)

	slog.Info("gRPC server is running", slog.String("address", address))
	grpcServer.Serve(listen)
}

func Shutdown(ctx context.Context, c *closer.Closer) {
	grpcServer.Stop()
	slog.Info("gRPC server stopped")
	c.Done()
}
