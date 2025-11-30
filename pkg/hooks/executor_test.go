package hooks

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestExportContextToEnv(t *testing.T) {
	ctx := ExportContext{
		ExportPath:   "/tmp/export.md",
		ExportFormat: "markdown",
		IssueCount:   42,
		Timestamp:    time.Date(2025, 11, 30, 10, 30, 0, 0, time.UTC),
	}

	env := ctx.ToEnv()

	expected := map[string]string{
		"BV_EXPORT_PATH":   "/tmp/export.md",
		"BV_EXPORT_FORMAT": "markdown",
		"BV_ISSUE_COUNT":   "42",
		"BV_TIMESTAMP":     "2025-11-30T10:30:00Z",
	}

	for _, e := range env {
		found := false
		for key, val := range expected {
			if e == key+"="+val {
				found = true
				break
			}
		}
		if !found {
			// Check if it's one of our expected keys
			for key := range expected {
				if len(e) > len(key) && e[:len(key)+1] == key+"=" {
					t.Errorf("unexpected value for %s: got %s", key, e)
				}
			}
		}
	}
}

func TestLoaderNoConfig(t *testing.T) {
	// Create a temp directory without hooks.yaml
	tmpDir := t.TempDir()

	loader := NewLoader(WithProjectDir(tmpDir))
	err := loader.Load()
	if err != nil {
		t.Fatalf("expected no error for missing config, got: %v", err)
	}

	if loader.HasHooks() {
		t.Error("expected no hooks when config is missing")
	}
}

func TestLoaderWithValidConfig(t *testing.T) {
	// Create temp directory with .bv/hooks.yaml
	tmpDir := t.TempDir()
	bvDir := filepath.Join(tmpDir, ".bv")
	if err := os.MkdirAll(bvDir, 0755); err != nil {
		t.Fatalf("failed to create .bv dir: %v", err)
	}

	configContent := `
hooks:
  pre-export:
    - name: validate
      command: echo "validating"
      timeout: 5s
  post-export:
    - name: notify
      command: echo "done"
      timeout: 10s
      env:
        CUSTOM_VAR: custom_value
`
	configPath := filepath.Join(bvDir, "hooks.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	loader := NewLoader(WithProjectDir(tmpDir))
	err := loader.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if !loader.HasHooks() {
		t.Error("expected hooks to be loaded")
	}

	preHooks := loader.GetHooks(PreExport)
	if len(preHooks) != 1 {
		t.Fatalf("expected 1 pre-export hook, got %d", len(preHooks))
	}
	if preHooks[0].Name != "validate" {
		t.Errorf("expected hook name 'validate', got %s", preHooks[0].Name)
	}
	if preHooks[0].Timeout != 5*time.Second {
		t.Errorf("expected timeout 5s, got %v", preHooks[0].Timeout)
	}
	if preHooks[0].OnError != "fail" {
		t.Errorf("expected on_error 'fail' for pre-export, got %s", preHooks[0].OnError)
	}

	postHooks := loader.GetHooks(PostExport)
	if len(postHooks) != 1 {
		t.Fatalf("expected 1 post-export hook, got %d", len(postHooks))
	}
	if postHooks[0].Name != "notify" {
		t.Errorf("expected hook name 'notify', got %s", postHooks[0].Name)
	}
	if postHooks[0].OnError != "continue" {
		t.Errorf("expected on_error 'continue' for post-export, got %s", postHooks[0].OnError)
	}
	if postHooks[0].Env["CUSTOM_VAR"] != "custom_value" {
		t.Errorf("expected CUSTOM_VAR env, got %v", postHooks[0].Env)
	}
}

func TestExecutorRunSimpleHook(t *testing.T) {
	config := &Config{
		Hooks: HooksByPhase{
			PreExport: []Hook{
				{
					Name:    "echo-test",
					Command: "echo hello",
					Timeout: 5 * time.Second,
					OnError: "fail",
				},
			},
		},
	}

	ctx := ExportContext{
		ExportPath:   "/tmp/test.md",
		ExportFormat: "markdown",
		IssueCount:   10,
		Timestamp:    time.Now(),
	}

	executor := NewExecutor(config, ctx)
	err := executor.RunPreExport()
	if err != nil {
		t.Fatalf("expected hook to succeed, got: %v", err)
	}

	results := executor.Results()
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if !results[0].Success {
		t.Errorf("expected success, got failure: %v", results[0].Error)
	}

	if results[0].Stdout != "hello" {
		t.Errorf("expected stdout 'hello', got %q", results[0].Stdout)
	}
}

func TestExecutorHookFailure(t *testing.T) {
	config := &Config{
		Hooks: HooksByPhase{
			PreExport: []Hook{
				{
					Name:    "fail-hook",
					Command: "exit 1",
					Timeout: 5 * time.Second,
					OnError: "fail",
				},
			},
		},
	}

	ctx := ExportContext{
		ExportPath:   "/tmp/test.md",
		ExportFormat: "markdown",
		IssueCount:   10,
		Timestamp:    time.Now(),
	}

	executor := NewExecutor(config, ctx)
	err := executor.RunPreExport()
	if err == nil {
		t.Error("expected error for failing pre-export hook")
	}
}

func TestExecutorHookFailureContinue(t *testing.T) {
	config := &Config{
		Hooks: HooksByPhase{
			PostExport: []Hook{
				{
					Name:    "fail-continue",
					Command: "exit 1",
					Timeout: 5 * time.Second,
					OnError: "continue",
				},
				{
					Name:    "should-run",
					Command: "echo still-running",
					Timeout: 5 * time.Second,
					OnError: "continue",
				},
			},
		},
	}

	ctx := ExportContext{
		ExportPath:   "/tmp/test.md",
		ExportFormat: "markdown",
		IssueCount:   10,
		Timestamp:    time.Now(),
	}

	executor := NewExecutor(config, ctx)
	err := executor.RunPostExport()
	if err != nil {
		t.Errorf("expected no error with on_error=continue, got: %v", err)
	}

	results := executor.Results()
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	// First hook should fail
	if results[0].Success {
		t.Error("expected first hook to fail")
	}

	// Second hook should succeed
	if !results[1].Success {
		t.Errorf("expected second hook to succeed, got: %v", results[1].Error)
	}
	if results[1].Stdout != "still-running" {
		t.Errorf("expected stdout 'still-running', got %q", results[1].Stdout)
	}
}

func TestExecutorHookTimeout(t *testing.T) {
	config := &Config{
		Hooks: HooksByPhase{
			PreExport: []Hook{
				{
					Name:    "slow-hook",
					Command: "sleep 10",
					Timeout: 100 * time.Millisecond,
					OnError: "fail",
				},
			},
		},
	}

	ctx := ExportContext{
		ExportPath:   "/tmp/test.md",
		ExportFormat: "markdown",
		IssueCount:   10,
		Timestamp:    time.Now(),
	}

	executor := NewExecutor(config, ctx)
	err := executor.RunPreExport()
	if err == nil {
		t.Error("expected timeout error")
	}

	results := executor.Results()
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if results[0].Success {
		t.Error("expected hook to fail due to timeout")
	}

	if results[0].Duration < 100*time.Millisecond {
		t.Errorf("expected duration >= 100ms, got %v", results[0].Duration)
	}
}

func TestExecutorEnvironmentVariables(t *testing.T) {
	config := &Config{
		Hooks: HooksByPhase{
			PreExport: []Hook{
				{
					Name:    "env-test",
					Command: "echo $BV_EXPORT_PATH $BV_ISSUE_COUNT",
					Timeout: 5 * time.Second,
					OnError: "fail",
				},
			},
		},
	}

	ctx := ExportContext{
		ExportPath:   "/custom/path.md",
		ExportFormat: "markdown",
		IssueCount:   99,
		Timestamp:    time.Now(),
	}

	executor := NewExecutor(config, ctx)
	err := executor.RunPreExport()
	if err != nil {
		t.Fatalf("expected success, got: %v", err)
	}

	results := executor.Results()
	if results[0].Stdout != "/custom/path.md 99" {
		t.Errorf("expected env vars in output, got %q", results[0].Stdout)
	}
}

func TestExecutorCustomEnvExpansion(t *testing.T) {
	// Set an env var to be expanded
	os.Setenv("TEST_HOOK_VAR", "expanded_value")
	defer os.Unsetenv("TEST_HOOK_VAR")

	config := &Config{
		Hooks: HooksByPhase{
			PreExport: []Hook{
				{
					Name:    "env-expand",
					Command: "echo $CUSTOM_VAR",
					Timeout: 5 * time.Second,
					OnError: "fail",
					Env: map[string]string{
						"CUSTOM_VAR": "${TEST_HOOK_VAR}",
					},
				},
			},
		},
	}

	ctx := ExportContext{
		ExportPath:   "/tmp/test.md",
		ExportFormat: "markdown",
		IssueCount:   10,
		Timestamp:    time.Now(),
	}

	executor := NewExecutor(config, ctx)
	err := executor.RunPreExport()
	if err != nil {
		t.Fatalf("expected success, got: %v", err)
	}

	results := executor.Results()
	if results[0].Stdout != "expanded_value" {
		t.Errorf("expected env expansion, got %q", results[0].Stdout)
	}
}

func TestExecutorSummary(t *testing.T) {
	config := &Config{
		Hooks: HooksByPhase{
			PreExport: []Hook{
				{
					Name:    "success-hook",
					Command: "echo ok",
					Timeout: 5 * time.Second,
					OnError: "continue",
				},
			},
			PostExport: []Hook{
				{
					Name:    "fail-hook",
					Command: "exit 1",
					Timeout: 5 * time.Second,
					OnError: "continue",
				},
			},
		},
	}

	ctx := ExportContext{
		ExportPath:   "/tmp/test.md",
		ExportFormat: "markdown",
		IssueCount:   10,
		Timestamp:    time.Now(),
	}

	executor := NewExecutor(config, ctx)
	_ = executor.RunPreExport()
	_ = executor.RunPostExport()

	summary := executor.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}

	// Should mention both success and failure
	if !contains(summary, "1 succeeded") || !contains(summary, "1 failed") {
		t.Errorf("summary should mention success and failure count: %s", summary)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestLoaderInvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	bvDir := filepath.Join(tmpDir, ".bv")
	if err := os.MkdirAll(bvDir, 0755); err != nil {
		t.Fatalf("failed to create .bv dir: %v", err)
	}

	// Invalid YAML
	configContent := `
hooks:
  pre-export:
    - name: [invalid yaml
`
	configPath := filepath.Join(bvDir, "hooks.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	loader := NewLoader(WithProjectDir(tmpDir))
	err := loader.Load()
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestLoaderSkipsEmptyCommands(t *testing.T) {
	tmpDir := t.TempDir()
	bvDir := filepath.Join(tmpDir, ".bv")
	if err := os.MkdirAll(bvDir, 0755); err != nil {
		t.Fatalf("failed to create .bv dir: %v", err)
	}

	configContent := `
hooks:
  pre-export:
    - name: empty
      command: ""
  post-export:
    - command: "   "
`
	configPath := filepath.Join(bvDir, "hooks.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	loader := NewLoader(WithProjectDir(tmpDir))
	if err := loader.Load(); err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if loader.HasHooks() {
		t.Error("expected hooks with empty commands to be skipped resulting in no hooks")
	}

	ws := loader.Warnings()
	if len(ws) == 0 {
		t.Fatal("expected warning for empty command")
	}
}
