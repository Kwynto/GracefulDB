package closer

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

const (
	CLOSER_TESTING_ITER_MIN int = 5
	CLOSER_TESTING_ITER     int = 100
)

func Test_New(t *testing.T) {
	res := New()
	if reflect.TypeOf(res) != reflect.TypeOf(&Closer{}) {
		t.Error("New() error = The function returns the wrong type")
	}
}

func Test_AddMsg(t *testing.T) {
	res := New()

	type args struct {
		msg string
	}

	tests := []struct {
		name string
		args args
		// index  int
		result string
	}{
		{
			name: "The correct result.",
			args: args{
				msg: "The correct result",
			},
			// index:  0,
			result: "[!] The correct result",
		},
		{
			name: "The result is also correct.",
			args: args{
				msg: "The result is also correct",
			},
			// index:  1,
			result: "[!] The result is also correct",
		},
	}

	for itt, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res.AddMsg(tt.args.msg)
			if res.Msgs[itt] != tt.result {
				t.Errorf("AddMsg() error: %v != %v", res.Msgs[itt], tt.result)
			}
		})
	}
}

func Test_Done(t *testing.T) {
	res := New()

	rand.Seed(time.Now().Unix())

	for i := 0; i < CLOSER_TESTING_ITER_MIN; i++ {
		t.Run("Done() function testing", func(t *testing.T) {
			randI := rand.Intn(75) + 5
			res.Counter = randI
			expectedResult := randI - 1
			res.Done()

			if res.Counter != expectedResult {
				t.Errorf("Done() error: %v != %v", res.Counter, expectedResult)
			}
		})
	}
}

func Test_AddHandler(t *testing.T) {
	res := New()

	hand := func(ctx context.Context, c *Closer) {
		c.AddMsg("true")
	}

	for i := 0; i < CLOSER_TESTING_ITER_MIN; i++ {
		t.Run("AddHandler() function testing", func(t *testing.T) {
			c := res.Counter
			res.AddHandler(hand)
			if res.Counter == c {
				t.Error("AddHandler() error: the counter is not working.")
			}

			res.funcs[fmt.Sprint(Handler(hand))](context.Background(), res)
			if res.Msgs[i] != "[!] true" {
				t.Error("AddHandler() error: incorrect handler execution.")
			}
		})
	}
}

func Test_DelHandler(t *testing.T) {
	res := New()

	hand := func(ctx context.Context, c *Closer) {
		c.AddMsg("true")
	}

	for i := 0; i < CLOSER_TESTING_ITER_MIN; i++ {
		t.Run("DelHandler() function testing", func(t *testing.T) {
			c := res.Counter
			res.AddHandler(hand)
			res.DelHandler(hand)

			if res.Counter != c {
				t.Error("DelHandler() error: the counter has not been reduced.")
			}

			if _, ok := res.funcs[fmt.Sprint(Handler(hand))]; ok {
				t.Error("DelHandler() error: the handler has not been deleted.")
			}
		})
	}
}

func Test_RunAndDelHandler(t *testing.T) {
	res := New()

	hand := func(ctx context.Context, c *Closer) {
		c.AddMsg("true")
		c.Done()
	}

	for i := 0; i < CLOSER_TESTING_ITER_MIN; i++ {
		t.Run("RunAndDelHandler() function testing", func(t *testing.T) {
			c := res.Counter
			res.AddHandler(hand)
			res.RunAndDelHandler(hand)

			time.Sleep(50 * time.Millisecond)

			if res.Counter != c {
				t.Errorf("RunAndDelHandler() error: the counter has not been reduced. %v != %v", res.Counter, c)
			}

			if _, ok := res.funcs[fmt.Sprint(Handler(hand))]; ok {
				t.Error("RunAndDelHandler() error: the handler has not been deleted.")
			}
		})
	}
}

func Test_Close(t *testing.T) {
	res_with_msg := New()
	res_without_msg := New()
	res_with_timout := New()

	hand_with_msg := func(ctx context.Context, c *Closer) {
		c.AddMsg("true")
		c.Done()
	}

	hand_without_msg := func(ctx context.Context, c *Closer) {
		c.Done()
	}

	hand_with_timeout := func(ctx context.Context, c *Closer) {
		time.Sleep(2 * time.Second)
		c.Done()
	}

	res_with_msg.AddHandler(hand_with_msg)
	res_without_msg.AddHandler(hand_without_msg)
	res_with_timout.AddHandler(hand_with_timeout)

	t.Run("Close() function testing with errors", func(t *testing.T) {

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := res_with_msg.Close(shutdownCtx)
		if err == nil {
			t.Error("Close() error: the error messages were not returned")
		}

		if res_with_msg.Counter != 0 {
			t.Error("Close() error: the counter has not been reduced.")
		}
	})

	t.Run("Close() function testing without errors", func(t *testing.T) {

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := res_without_msg.Close(shutdownCtx)
		if err != nil {
			t.Error("Close() error: the error messages were returned")
		}

		if res_without_msg.Counter != 0 {
			t.Error("Close() error: the counter has not been reduced.")
		}
	})

	t.Run("Close() function testing with timeout error", func(t *testing.T) {

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := res_with_timout.Close(shutdownCtx)
		if err == nil {
			t.Error("Close() error: the context error were not returned")
		}
	})
}

func Test_AddHandler_Default(t *testing.T) {
	CloseProcs = New()

	hand := func(ctx context.Context, c *Closer) {
		c.AddMsg("true")
	}

	for i := 0; i < CLOSER_TESTING_ITER_MIN; i++ {
		t.Run("AddHandler() function testing (default)", func(t *testing.T) {
			c := CloseProcs.Counter
			AddHandler(hand)
			if CloseProcs.Counter == c {
				t.Error("AddHandler() error: the counter is not working.")
			}

			CloseProcs.funcs[fmt.Sprint(Handler(hand))](context.Background(), CloseProcs)
			if CloseProcs.Msgs[i] != "[!] true" {
				t.Error("AddHandler() error: incorrect handler execution.")
			}
		})
	}
}

func Test_DelHandler_Default(t *testing.T) {
	CloseProcs = New()

	hand := func(ctx context.Context, c *Closer) {
		c.AddMsg("true")
	}

	for i := 0; i < CLOSER_TESTING_ITER_MIN; i++ {
		t.Run("DelHandler() function testing (default)", func(t *testing.T) {
			c := CloseProcs.Counter
			AddHandler(hand)
			DelHandler(hand)

			if CloseProcs.Counter != c {
				t.Error("DelHandler() error: the counter has not been reduced.")
			}

			if _, ok := CloseProcs.funcs[fmt.Sprint(Handler(hand))]; ok {
				t.Error("DelHandler() error: the handler has not been deleted.")
			}
		})
	}
}

func Test_RunAndDelHandler_Default(t *testing.T) {
	CloseProcs = New()

	hand := func(ctx context.Context, c *Closer) {
		c.AddMsg("true")
		c.Done()
	}

	for i := 0; i < CLOSER_TESTING_ITER_MIN; i++ {
		t.Run("RunAndDelHandler() function testing (default)", func(t *testing.T) {
			c := CloseProcs.Counter
			AddHandler(hand)
			RunAndDelHandler(hand)

			time.Sleep(50 * time.Millisecond)

			if CloseProcs.Counter != c {
				t.Errorf("RunAndDelHandler() error: the counter has not been reduced. %v != %v", CloseProcs.Counter, c)
			}

			if _, ok := CloseProcs.funcs[fmt.Sprint(Handler(hand))]; ok {
				t.Error("RunAndDelHandler() error: the handler has not been deleted.")
			}
		})
	}
}

func Test_Close_Default(t *testing.T) {
	hand_with_msg := func(ctx context.Context, c *Closer) {
		c.AddMsg("true")
		c.Done()
	}

	hand_without_msg := func(ctx context.Context, c *Closer) {
		c.Done()
	}

	hand_with_timeout := func(ctx context.Context, c *Closer) {
		time.Sleep(2 * time.Second)
		c.Done()
	}

	CloseProcs = New()
	AddHandler(hand_with_msg)

	t.Run("Close() function testing with errors (default)", func(t *testing.T) {

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := Close(shutdownCtx)
		if err == nil {
			t.Error("Close() error: the error messages were not returned")
		}

		if CloseProcs.Counter != 0 {
			t.Error("Close() error: the counter has not been reduced.")
		}
	})

	CloseProcs = New()
	AddHandler(hand_without_msg)
	// res_with_timout.AddHandler(hand_with_timeout)

	t.Run("Close() function testing without errors (default)", func(t *testing.T) {

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := Close(shutdownCtx)
		if err != nil {
			t.Error("Close() error: the error messages were returned")
		}

		if CloseProcs.Counter != 0 {
			t.Error("Close() error: the counter has not been reduced.")
		}
	})

	CloseProcs = New()
	AddHandler(hand_with_timeout)

	t.Run("Close() function testing with timeout error (default)", func(t *testing.T) {

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := Close(shutdownCtx)
		if err == nil {
			t.Error("Close() error: the context error were not returned")
		}
	})
}
