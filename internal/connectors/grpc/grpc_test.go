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

func Test_Query(t *testing.T) {

	testingEntity := tMessageServer{}

	t.Run("Query() function testing", func(t *testing.T) {
		ctx := context.Background()
		req := gs.Request{
			Instruction: "instruction",
			Placeholder: []string{},
		}

		resp, err := testingEntity.Query(ctx, &req) // calling the tested function

		if err != nil {
			t.Error("Query() error.")
		}

		if reflect.TypeOf(resp) != reflect.TypeOf(&gs.Response{}) {
			t.Error("Query() error = The function returns the wrong type")
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
		if closer.StCloseProcs.Counter != 0 {
			t.Errorf("Start() error: %v.", closer.StCloseProcs.Counter)
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

		Shutdown(context.Background(), closer.StCloseProcs)

		if closer.StCloseProcs.Counter != 0 {
			t.Errorf("Shutdown() error: %v.", closer.StCloseProcs.Counter)
		}
	})
}
