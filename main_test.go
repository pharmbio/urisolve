package main

import (
	"testing"
)

func TestValidQuery(t *testing.T) {
	queries := map[string]bool{
		"? http://ex.org/ex; http://ex.org/ex": false,
		"? http://ex.org/;ex http://ex.org/ex": false,
		"":                                       false,
		"? ? ?":                                  true,
		"? ? http://ex.org/ex":                   true,
		"? http://ex.org/ex ?":                   true,
		"http://ex.org/ex ? ?":                   true,
		"http://ex.org/ex ? http://ex.org/ex":    true,
		"http://ex.org/ex http://ex.org/ex ?":    true,
		"? http://ex.org/ex http://ex.org/ex":    true,
		"? ? https://ex.org/foo-bar.php#foo_bar": true,
	}
	for q, shouldBeOk := range queries {
		if validQuery(q) != shouldBeOk {
			t.Errorf("Query was %v. Expected %v: %s", !shouldBeOk, shouldBeOk, q)
		}
	}
}
