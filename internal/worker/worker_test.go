package worker

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	switch os.Getenv("SCRYER_ENGINE_FAKE_JVM") {
	case "echo":
		fmt.Fprintln(os.Stderr, "worker progress")
		_, _ = io.Copy(os.Stdout, os.Stdin)
		os.Exit(0)
	case "wait":
		time.Sleep(30 * time.Second)
		os.Exit(0)
	case "exit":
		os.Exit(7)
	}
	os.Exit(m.Run())
}
func Test_should_keep_json_separate_when_forwarding_worker_streams(t *testing.T) {
	t.Setenv("SCRYER_ENGINE_FAKE_JVM", "echo")
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	var out, diagnostics bytes.Buffer
	request := `{"requestId":"opaque","operation":"scan","future":null}`
	if err := run(context.Background(), executable, "unused", strings.NewReader(request), &out, &diagnostics); err != nil {
		t.Fatal(err)
	}
	if out.String() != request || !strings.Contains(diagnostics.String(), "worker progress") {
		t.Fatal("worker streams were changed or mixed")
	}
}
func Test_should_stop_worker_when_context_is_cancelled(t *testing.T) {
	t.Setenv("SCRYER_ENGINE_FAKE_JVM", "wait")
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	started := time.Now()
	if err := run(ctx, executable, "unused", nil, io.Discard, io.Discard); err != context.DeadlineExceeded {
		t.Fatalf("unexpected cancellation: %v", err)
	}
	if time.Since(started) > 5*time.Second {
		t.Fatal("worker did not stop")
	}
}
func Test_should_return_error_when_worker_exits_unsuccessfully(t *testing.T) {
	t.Setenv("SCRYER_ENGINE_FAKE_JVM", "exit")
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	if err := run(context.Background(), executable, "unused", nil, io.Discard, io.Discard); err == nil {
		t.Fatal("worker failure was ignored")
	}
}
func Test_should_select_engine_runtime_when_target_java_is_different(t *testing.T) {
	installation := filepath.Join(t.TempDir(), "bin", "scryer")
	bundled := filepath.Join(filepath.Dir(installation), "..", "runtime", "bin", "java")
	if got := javaExecutable(installation, "", "/target-java8"); got != "/target-java8/bin/java" {
		t.Fatal(got)
	}
	if err := os.MkdirAll(filepath.Dir(bundled), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(bundled, []byte("runtime"), 0755); err != nil {
		t.Fatal(err)
	}
	if got := javaExecutable(installation, "", "/target-java8"); got != bundled {
		t.Fatal(got)
	}
	if got := javaExecutable(installation, "/explicit-java21", "/target-java8"); got != "/explicit-java21/bin/java" {
		t.Fatal(got)
	}
}
