package config

import (
	"testing"
)

func TestParse(t *testing.T) {
	expected := "go-gettingstarted"
	cfg, err := ParseManual()
	if err != nil {
		t.Fatalf("unexpected error when calling ParseManual(): %s", err)
	}

	if cfg.AppName != expected {
		t.Errorf("expected: %s, got: %s", expected, cfg.AppName)
	}
}

func TestParseLib(t *testing.T) {
	want := "go-gettingstarted"
	cfg, err := Parse()
	if err != nil {
		t.Fatalf("unexpected error when calling Parse(): %s", err)
	}

	if cfg.AppName != want {
		t.Errorf("want: %s, got: %s", want, cfg.AppName)
	}
}
