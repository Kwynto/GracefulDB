package sqlanalyzer

import "testing"

func Test_Request(t *testing.T) {
	t.Run("Request() function testing", func(t *testing.T) {
		ticket := ""
		instruction := ""
		placeholder := []string{}
		res := Request(&ticket, &instruction, &placeholder)

		if *res == "" {
			t.Errorf("Request() error: empty result.")
		}
	})

}
