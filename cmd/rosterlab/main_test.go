package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestCLIRejectsMalformedInputWithoutPlan(t *testing.T) {
	for _, input := range []string{`{"topic":"api","unknown":true}`, `{"topic":"api"} {}`, `null`, strings.Repeat(" ", (1<<20)+1)} {
		var out bytes.Buffer
		if err := run(nil, strings.NewReader(input), &out); err == nil || out.Len() != 0 {
			t.Fatalf("err %v, output %q", err, out.String())
		}
	}
}

func TestCLIProducesSortedPlan(t *testing.T) {
	var out bytes.Buffer
	input := `{"topic":"api","current":[{"id":"team","email":"team@example.invalid","topics":["api"]}],"required":["team"]}`
	if err := run(nil, strings.NewReader(input), &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), `"team@example.invalid"`) {
		t.Fatal(out.String())
	}
}
