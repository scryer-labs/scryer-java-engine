// Package cli exposes the initial mock engine's process interface.
package cli

import (
	"encoding/json"
	"fmt"
	"io"
)

type Identity struct {
	ProtocolVersion int    `json:"protocolVersion"`
	ID              string `json:"id"`
	Name            string `json:"name"`
	Version         string `json:"version"`
	Source          string `json:"source"`
	Mock            bool   `json:"mock"`
}

func Run(args []string, out, diagnostics io.Writer) int {
	if len(args) == 0 || (len(args) == 1 && (args[0] == "--help" || args[0] == "help")) {
		fmt.Fprintln(out, "Usage: scryer-java-engine describe\nMock Java engine; no scanning or installation is implemented.")
		return 0
	}
	if len(args) != 1 || args[0] != "describe" {
		fmt.Fprintln(diagnostics, "scryer-java-engine: expected describe")
		return 2
	}
	identity := Identity{
		ProtocolVersion: 1,
		ID:              "java",
		Name:            "Scryer Java Engine",
		Version:         "0.0.0-dev",
		Source:          "github.com/scryer-labs/scryer-java-engine",
		Mock:            true,
	}
	if err := json.NewEncoder(out).Encode(identity); err != nil {
		fmt.Fprintln(diagnostics, "scryer-java-engine:", err)
		return 1
	}
	return 0
}
