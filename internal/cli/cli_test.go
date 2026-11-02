package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"testing"
)

func Test_should_return_one_engine_identity_when_describe_is_requested(t *testing.T) {
	var out, diagnostics bytes.Buffer
	if status := Run(context.Background(), []string{"describe"}, nil, &out, &diagnostics); status != 0 {
		t.Fatalf("unexpected exit status: %d", status)
	}
	decoder := json.NewDecoder(&out)
	var identity Identity
	if err := decoder.Decode(&identity); err != nil {
		t.Fatal(err)
	}
	if identity.ProtocolVersion != 1 || identity.ID != "java" || identity.Mock || identity.Source != "github.com/scryer-labs/scryer-java-engine" {
		t.Fatalf("unexpected identity: %+v", identity)
	}
	if decoder.Decode(new(any)) != io.EOF || diagnostics.Len() != 0 {
		t.Fatal("describe must return exactly one JSON result without diagnostics")
	}
}

func Test_should_reject_unknown_operations_when_operation_is_not_supported(t *testing.T) {
	for _, args := range [][]string{{"scan"}, {"install"}, {"describe", "extra"}} {
		var out, diagnostics bytes.Buffer
		if Run(context.Background(), args, nil, &out, &diagnostics) != 2 || out.Len() != 0 || diagnostics.Len() == 0 {
			t.Fatalf("unexpected result for %v", args)
		}
	}
}

func Test_should_describe_relocatable_package_when_manifest_is_requested(t *testing.T) {
	var out, diagnostics bytes.Buffer
	if status := Run(context.Background(), []string{"package-manifest"}, nil, &out, &diagnostics); status != 0 {
		t.Fatal(status, diagnostics.String())
	}
	var manifest struct {
		SchemaVersion int      `json:"schemaVersion"`
		Engine        Identity `json:"engine"`
		Executable    string   `json:"executable"`
		OS            string   `json:"os"`
		Arch          string   `json:"arch"`
	}
	if err := json.Unmarshal(out.Bytes(), &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.SchemaVersion != 1 || manifest.Engine.ID != "java" || manifest.Engine.Mock || manifest.Executable != "bin/scryer-java-engine" || manifest.OS == "" || manifest.Arch == "" {
		t.Fatalf("invalid package manifest: %+v", manifest)
	}
}
