package cli

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"
)

func Test_should_return_one_mock_identity_when_describe_is_requested(t *testing.T) {
	var out, diagnostics bytes.Buffer
	if status := Run([]string{"describe"}, &out, &diagnostics); status != 0 {
		t.Fatalf("unexpected exit status: %d", status)
	}
	decoder := json.NewDecoder(&out)
	var identity Identity
	if err := decoder.Decode(&identity); err != nil {
		t.Fatal(err)
	}
	if identity.ProtocolVersion != 1 || identity.ID != "java" || !identity.Mock || identity.Source != "github.com/scryer-labs/scryer-java-engine" {
		t.Fatalf("unexpected identity: %+v", identity)
	}
	if decoder.Decode(new(any)) != io.EOF || diagnostics.Len() != 0 {
		t.Fatal("describe must return exactly one JSON result without diagnostics")
	}
}

func Test_should_reject_unknown_operations_when_engine_is_only_a_mock(t *testing.T) {
	for _, args := range [][]string{{"scan"}, {"install"}, {"describe", "extra"}} {
		var out, diagnostics bytes.Buffer
		if Run(args, &out, &diagnostics) != 2 || out.Len() != 0 || diagnostics.Len() == 0 {
			t.Fatalf("unexpected result for %v", args)
		}
	}
}
