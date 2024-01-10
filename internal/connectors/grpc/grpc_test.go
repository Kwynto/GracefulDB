package grpc

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/Kwynto/GracefulDB/internal/config"
	gs "github.com/Kwynto/GracefulDB/internal/connectors/grpc/proto/graceful_service"
	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
)

func Test_SQuery(t *testing.T) {

	testingEntity := tMessageServer{}

	t.Run("SQuery() function testing", func(t *testing.T) {
		ctx := context.Background()
		req := gs.SRequest{
			Instruction: "instruction",
			Placeholder: []string{},
		}

		resp, err := testingEntity.SQuery(ctx, &req) // calling the tested function

		if err != nil {
			t.Error("SQuery() error.")
		}

		if reflect.TypeOf(resp) != reflect.TypeOf(&gs.SResponse{}) {
			t.Error("SQuery() error = The function returns the wrong type")
		}
	})
}

func Test_VQuery(t *testing.T) {

	testingEntity := tMessageServer{}

	t.Run("VQuery() function testing", func(t *testing.T) {
		ctx := context.Background()
		req := gs.VRequest{
			Instruction: "instruction",
		}

		resp, err := testingEntity.VQuery(ctx, &req) // calling the tested function

		if err != nil {
			t.Error("VQuery() error.")
		}

		if reflect.TypeOf(resp) != reflect.TypeOf(&gs.VResponse{}) {
			t.Error("VQuery() error = The function returns the wrong type")
		}
	})
}

func Test_Start(t *testing.T) {
	t.Run("Start() function testing", func(t *testing.T) {
		tf := "../../../../config/develop.yaml"
		config.MustLoad(tf)
		config.DefaultConfig.GrpcConnector.Address = "256.256.256.256" // Creating an error

		go Start(&config.DefaultConfig)
		time.Sleep(1 * time.Second)

		config.MustLoad(tf)
		if closer.CloseProcs.Counter != 0 {
			t.Errorf("Start() error: %v.", closer.CloseProcs.Counter)
		}
	})
}

func Test_Shutdown(t *testing.T) {
	t.Run("Shutdown() function testing", func(t *testing.T) {
		tf := "../../../../config/develop.yaml"
		config.MustLoad(tf)

		go Start(&config.DefaultConfig)
		closer.AddHandler(Shutdown)
		time.Sleep(1 * time.Second)

		Shutdown(context.Background(), closer.CloseProcs)

		if closer.CloseProcs.Counter != 0 {
			t.Errorf("Shutdown() error: %v.", closer.CloseProcs.Counter)
		}
	})
}
