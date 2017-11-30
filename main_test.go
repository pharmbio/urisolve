package main

import (
	"testing"
)

func TestValidUri(t *testing.T) {
	uris := map[string]bool{
		"http://ex.org/ex;":                  false,
		"http://ex.org/;ex":                  false,
		"https://ex.org/foo-bar.php#foo_bar": true,
	}
	var explanation = map[bool]string{
		true:  "allowed",
		false: "forbidden",
	}
	for uri, shouldBeOk := range uris {
		if validUri(uri) != shouldBeOk {
			t.Errorf("Uri pattern was %s. Expected it to be %s: %s", explanation[!shouldBeOk], explanation[shouldBeOk], uri)
		}
	}
}
