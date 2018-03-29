package main

import "testing"

func TestFormatYaml(t *testing.T) {
	_, err := formatYaml([]byte(`illegal content`))
	if err == nil {
		t.Fatal("err expected")
	}
}
