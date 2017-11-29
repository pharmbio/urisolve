package main

import (
	"testing"
)

func TestValidQuery(t *testing.T) {
	saneQuery := "? http://ex.org/ex http://ex.org/ex"
	saneOk := validQuery(saneQuery)
	if !saneOk {
		t.Errorf("Sane query returned error: %s", saneQuery)
	}

	nastyQuery := "? http://ex.org/ex; http://ex.org/ex"
	nastyOk := validQuery(nastyQuery)
	if nastyOk {
		t.Errorf("Nasty query not stopped: %s", nastyQuery)
	}
}
