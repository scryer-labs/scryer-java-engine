// Package worker owns JVM discovery and execution inside the Java engine.
package worker

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ExecuteAt resolves runtime and jars relative to the launcher or shared library.
func ExecuteAt(ctx context.Context, installation string, in io.Reader, out, diagnostics io.Writer, ownGroup bool) error {
	installation, err := filepath.EvalSymlinks(installation)
	if err != nil {
		return err
	}
	classpath := os.Getenv("SCRYER_JAVA_CLASSPATH")
	if classpath == "" {
		jars, err := filepath.Glob(filepath.Join(filepath.Dir(installation), "..", "lib", "java", "*.jar"))
		if err != nil {
			return err
		}
		if len(jars) == 0 {
			return fmt.Errorf("Java worker is not packaged: run the engine's scripts/build or set SCRYER_JAVA_CLASSPATH")
		}
		classpath = strings.Join(jars, string(os.PathListSeparator))
	}
	executable := javaExecutable(installation, os.Getenv("SCRYER_ANALYZER_JAVA_HOME"), os.Getenv("JAVA_HOME"))
	return runWorker(ctx, executable, classpath, in, out, diagnostics, ownGroup)
}

func run(ctx context.Context, executable, classpath string, in io.Reader, out, diagnostics io.Writer) error {
	return runWorker(ctx, executable, classpath, in, out, diagnostics, false)
}

func runWorker(ctx context.Context, executable, classpath string, in io.Reader, out, diagnostics io.Writer, ownGroup bool) error {
	command := exec.CommandContext(ctx, executable, "-cp", classpath, "com.edwardnoaland.scryer.worker.WorkerKt")
	// The ABI worker owns a separate group. The executable wrapper inherits Core’s group.
	if ownGroup {
		configureWorkerGroup(command)
	}
	command.Stdin = in
	command.Stdout = out
	command.Stderr = diagnostics
	command.WaitDelay = 3 * time.Second
	if err := command.Run(); err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return fmt.Errorf("Java worker failed (requires Java 21; set SCRYER_ANALYZER_JAVA_HOME): %w", err)
	}
	return nil
}

// The analyzer runtime is separate from a target project's build/test toolchain.
func javaExecutable(installation, overrideHome, javaHome string) string {
	if overrideHome != "" {
		return filepath.Join(overrideHome, "bin", "java")
	}
	bundled := filepath.Join(filepath.Dir(installation), "..", "runtime", "bin", "java")
	if info, err := os.Stat(bundled); err == nil && !info.IsDir() {
		return bundled
	}
	if javaHome != "" {
		return filepath.Join(javaHome, "bin", "java")
	}
	return "java"
}
