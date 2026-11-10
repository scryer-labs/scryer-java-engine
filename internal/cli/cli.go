// Package cli exposes the Java engine process interface.
package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/scryer-labs/scryer-java-engine/internal/engine"
	"github.com/spf13/cobra"
)

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
			return json.NewEncoder(cmd.OutOrStdout()).Encode(engine.Describe())
		}),
		command("package-manifest", "Print the engine package manifest as JSON", func(cmd *cobra.Command) error {
			return json.NewEncoder(cmd.OutOrStdout()).Encode(engine.PackageManifest())
		}),
		command("execute", "Read one JSON request from stdin and print the engine response", func(cmd *cobra.Command) error {
			request, err := engine.ReadRequest(cmd.InOrStdin())
			if err != nil {
				return err
			}
			installation, err := os.Executable()
			if err != nil {
				return err
			}
			return request.Execute(cmd.Context(), installation, cmd.OutOrStdout(), cmd.ErrOrStderr(), false)
		}),
	)
	if err := root.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(root.ErrOrStderr(), "scryer-java-engine:", err)
		return exitCode
	}
	return 0
}
