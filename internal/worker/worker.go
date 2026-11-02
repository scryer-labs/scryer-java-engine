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

func Execute(ctx context.Context, in io.Reader, out, diagnostics io.Writer) error {
	installation, err := os.Executable()
	if err != nil {
		return err
	}
	installation, err = filepath.EvalSymlinks(installation)
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
	return run(ctx, executable, classpath, in, out, diagnostics)
}

func run(ctx context.Context, executable, classpath string, in io.Reader, out, diagnostics io.Writer) error {
	command := exec.CommandContext(ctx, executable, "-cp", classpath, "com.edwardnoaland.scryer.worker.WorkerKt")
	// Inherit the wrapper's process group so Core cancellation reaches both processes.
	// Starting another group here would let the JVM escape Core's group cancellation.
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
