package sqlanalyzer

import "testing"

func Test_Request(t *testing.T) {
	t.Run("Request() function testing", func(t *testing.T) {
		instruction := ""
		placeholder := []string{}
		res := Request(&instruction, &placeholder)

		if *res == "" {
			t.Errorf("Request() error: empty result.")
		}
	})

}
