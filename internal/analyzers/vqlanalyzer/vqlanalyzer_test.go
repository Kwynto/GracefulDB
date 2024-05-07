package vqlanalyzer

import "testing"

func Test_Request(t *testing.T) {
	t.Run("Request() function testing", func(t *testing.T) {
		sTicket := ""
		sInstruction := ""
		slPlaceholder := []string{}
		sResult := Request(sTicket, sInstruction, slPlaceholder)

		if sResult == "" {
			t.Errorf("Request() error: empty result.")
		}
	})

}
