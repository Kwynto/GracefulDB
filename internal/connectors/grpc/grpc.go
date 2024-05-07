package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

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

var sAddress string

var stMessageServer tMessageServer
var stGrpcServer *grpc.Server

func (tMessageServer) Query(ctx context.Context, r *gs.Request) (stResponse *gs.Response, err error) {
	sOperation := "internal -> connectors -> gRPC -> SQuery"
	defer func() { e.Wrapper(sOperation, err) }()

	slog.Debug("Request received", slog.String("instruction", r.Instruction), slog.String("placeholder", fmt.Sprint(r.Placeholder)))

	stResponse = &gs.Response{
		Message: vqlanalyzer.Request(r.Ticket, r.Instruction, r.Placeholder),
	}
	slog.Debug("Response sent", slog.String("response", stResponse.Message))

	return stResponse, nil
}

func Start(cfg *config.TConfig) {
	sAddress = fmt.Sprintf("%s:%s", cfg.GrpcConnector.Address, cfg.GrpcConnector.Port)
	inListen, err := net.Listen("tcp", sAddress)
	if err != nil {
		slog.Error("Failed to start gRPC-listener", slog.String("err", err.Error()))
		return
	}

	stGrpcServer = grpc.NewServer()
	gs.RegisterGracefulServiceServer(stGrpcServer, stMessageServer)

	slog.Info("gRPC server is running", slog.String("address", sAddress))
	stGrpcServer.Serve(inListen)
}

func Shutdown(ctx context.Context, c *closer.TCloser) {
	stGrpcServer.Stop()
	slog.Info("gRPC server stopped")
	c.Done()
}
