// Package cli exposes the Java engine process interface.
package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"runtime"

	"github.com/scryer-labs/scryer-java-engine/internal/worker"
	"github.com/spf13/cobra"
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
	exitCode := 2
	root := &cobra.Command{
		Use:               "scryer-java-engine",
		Short:             "Run the standalone Java engine",
		Args:              cobra.NoArgs,
		RunE:              func(cmd *cobra.Command, _ []string) error { return cmd.Help() },
		SilenceErrors:     true,
		SilenceUsage:      true,
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}
	// A nil slice tells Cobra to read the host process arguments.
	root.SetArgs(append([]string{}, args...))
	root.SetIn(in)
	root.SetOut(out)
	root.SetErr(diagnostics)

	identity := Identity{
		ProtocolVersion: 1,
		ID:              "java",
		Name:            "Scryer Java Engine",
		Version:         "0.0.0-dev",
		Source:          "github.com/scryer-labs/scryer-java-engine",
	}
	capabilities := struct {
		Languages  []string `json:"languages"`
		Operations []string `json:"operations"`
		ScanFacts  []string `json:"scanFacts"`
	}{[]string{"Java"}, []string{"scan", "analyze"}, []string{"dependencies"}}
	manifest := struct {
		Capabilities  any      `json:"capabilities"`
		SchemaVersion int      `json:"schemaVersion"`
		Engine        Identity `json:"engine"`
		Library       string   `json:"library,omitempty"`
		Executable    string   `json:"executable"`
		OS            string   `json:"os"`
		Arch          string   `json:"arch"`
	}{capabilities, 1, identity, libraryPath(), "bin/scryer-java-engine", runtime.GOOS, runtime.GOARCH}

	command := func(name, description string, run func(*cobra.Command) error) *cobra.Command {
		return &cobra.Command{
			Use: name, Short: description, Args: cobra.NoArgs,
			RunE: func(cmd *cobra.Command, _ []string) error {
				exitCode = 1
				return run(cmd)
			},
		}
	}
	root.AddCommand(
		command("describe", "Print the engine identity as JSON", func(cmd *cobra.Command) error {
			return json.NewEncoder(cmd.OutOrStdout()).Encode(identity)
		}),
		command("package-manifest", "Print the engine package manifest as JSON", func(cmd *cobra.Command) error {
			return json.NewEncoder(cmd.OutOrStdout()).Encode(manifest)
		}),
		command("execute", "Read one JSON request from stdin and print the JVM response", func(cmd *cobra.Command) error {
			return worker.Execute(cmd.Context(), cmd.InOrStdin(), cmd.OutOrStdout(), cmd.ErrOrStderr())
		}),
	)
	if err := root.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(root.ErrOrStderr(), "scryer-java-engine:", err)
		return exitCode
	}
	return 0
}

func libraryPath() string {
	switch runtime.GOOS {
	case "darwin":
		return "bin/libscryer-java-engine.dylib"
	case "linux":
		return "bin/libscryer-java-engine.so"
	default:
		return ""
	}
}
