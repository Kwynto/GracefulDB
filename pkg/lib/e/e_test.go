package e

import (
	"fmt"
	"testing"
)

func TestWrapper(t *testing.T) {
	type args struct {
		msg string
		err error
	}

	stTests := []struct {
		name      string
		args      args
		isWantErr bool
	}{
		{
			name: "Error wrapping test, if error is not nil.",
			args: args{
				msg: "Msg",
				err: fmt.Errorf("%s", "error msg"),
			},
			isWantErr: true,
		},
		{
			name: "Error wrapping test, if error is nil.",
			args: args{
				msg: "Msg",
				err: nil,
			},
			isWantErr: false,
		},
		{
			name: "Error logging.",
			args: args{
				msg: "This is not an error, this diagnostic message is part of a test",
				err: fmt.Errorf("%s", "test msg"),
			},
			isWantErr: true,
		},
	}

	for _, tt := range stTests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Wrapper(tt.args.msg, tt.args.err); (err != nil) != tt.isWantErr {
				t.Errorf("WrapIfErr() error = %v, wantErr %v", err, tt.isWantErr)
			}
		})
	}
}

func Benchmark_Wrapper(b *testing.B) {
	errB := fmt.Errorf("%s", "error msg")
	for i := 0; i < b.N; i++ {
		_ = Wrapper("Msg", errB) // calling the tested function
	}
}
