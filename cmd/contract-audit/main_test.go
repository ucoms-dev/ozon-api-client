package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunAuditsProvidedSwaggerAndClientDirectory(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	clientDir := filepath.Join(tempDir, "client")
	if err := os.Mkdir(clientDir, 0o700); err != nil {
		t.Fatalf("mkdir client fixture: %v", err)
	}
	goSource := `package fixture
import "net/http"
func current() { Request(http.MethodPost, "/v1/current") }
`
	if err := os.WriteFile(filepath.Join(clientDir, "client.go"), []byte(goSource), 0o600); err != nil {
		t.Fatalf("write Go fixture: %v", err)
	}
	swaggerPath := filepath.Join(tempDir, "swagger.json")
	swagger := `{"openapi":"3.0.0","paths":{"/v1/current":{"post":{"operationId":"Current"}}}}`
	if err := os.WriteFile(swaggerPath, []byte(swagger), 0o600); err != nil {
		t.Fatalf("write Swagger fixture: %v", err)
	}

	var stdout bytes.Buffer
	if err := run([]string{"-swagger", swaggerPath, "-client", clientDir, "-date", "2026-09-02"}, &stdout); err != nil {
		t.Fatalf("run: %v", err)
	}
	output := stdout.String()
	for _, want := range []string{
		"- Client operations: 1",
		"- Swagger operations: 1",
		"- Exact method/path matches: 1",
		"- Swagger SHA-256: `",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("output does not contain %q:\n%s", want, output)
		}
	}
}

func TestRunWritesReportToOutputFile(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	clientDir := filepath.Join(tempDir, "client")
	if err := os.Mkdir(clientDir, 0o700); err != nil {
		t.Fatalf("mkdir client fixture: %v", err)
	}
	if err := os.WriteFile(filepath.Join(clientDir, "client.go"), []byte(`package fixture
import "net/http"
func current() { Request(http.MethodPost, "/v1/current") }
`), 0o600); err != nil {
		t.Fatalf("write Go fixture: %v", err)
	}
	swaggerPath := filepath.Join(tempDir, "swagger.json")
	if err := os.WriteFile(swaggerPath, []byte(`{"openapi":"3.0.0","paths":{"/v1/current":{"post":{"operationId":"Current"}}}}`), 0o600); err != nil {
		t.Fatalf("write Swagger fixture: %v", err)
	}
	reportPath := filepath.Join(tempDir, "report.md")

	var stdout bytes.Buffer
	if err := run([]string{"-swagger", swaggerPath, "-client", clientDir, "-output", reportPath}, &stdout); err != nil {
		t.Fatalf("run: %v", err)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	report, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("read report: %v", err)
	}
	if !strings.Contains(string(report), "- Exact method/path matches: 1") {
		t.Fatalf("unexpected report:\n%s", report)
	}
}
