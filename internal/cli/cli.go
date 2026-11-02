// Package cli exposes the Java engine process interface.
package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"runtime"

	"github.com/scryer-labs/scryer-java-engine/internal/worker"
)

type Identity struct {
	ProtocolVersion int    `json:"protocolVersion"`
	ID              string `json:"id"`
	Name            string `json:"name"`
	Version         string `json:"version"`
	Source          string `json:"source"`
	Mock            bool   `json:"mock"`
}

func Run(ctx context.Context, args []string, in io.Reader, out, diagnostics io.Writer) int {
	if len(args) == 0 || (len(args) == 1 && (args[0] == "--help" || args[0] == "help")) {
		fmt.Fprintln(out, "Usage: scryer-java-engine describe | execute\nexecute reads one v1 scan/analyze JSON request from stdin and writes the JVM response.")
		return 0
	}
	if len(args) == 1 && args[0] == "execute" {
		if err := worker.Execute(ctx, in, out, diagnostics); err != nil {
			fmt.Fprintln(diagnostics, "scryer-java-engine:", err)
			return 1
		}
		return 0
	}
	if len(args) != 1 || (args[0] != "describe" && args[0] != "package-manifest") {
		fmt.Fprintln(diagnostics, "scryer-java-engine: expected describe, package-manifest or execute")
		return 2
	}
	identity := Identity{
		ProtocolVersion: 1,
		ID:              "java",
		Name:            "Scryer Java Engine",
		Version:         "0.0.0-dev",
		Source:          "github.com/scryer-labs/scryer-java-engine",
		Mock:            false,
	}
	var response any = identity
	if args[0] == "package-manifest" {
		response = struct {
			SchemaVersion int      `json:"schemaVersion"`
			Engine        Identity `json:"engine"`
			Executable    string   `json:"executable"`
			OS            string   `json:"os"`
			Arch          string   `json:"arch"`
		}{1, identity, "bin/scryer-java-engine", runtime.GOOS, runtime.GOARCH}
	}
	if err := json.NewEncoder(out).Encode(response); err != nil {
		fmt.Fprintln(diagnostics, "scryer-java-engine:", err)
		return 1
	}
	return 0
}
