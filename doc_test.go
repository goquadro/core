package core

import "testing"

func TestsanitizeUrl(t *testing.T) {
	u := "http://www.goquadro.com"
	if parsed, err := sanitizeUrl(u); err != nil {
		t.Error("Address not correctly parsed")
	}
}
