package contractaudit

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAuditClassifiesExactDeprecatedMissingAndMethodMismatchOperations(t *testing.T) {
	t.Parallel()

	sourceDir := t.TempDir()
	source := `package fixture
import (
	"context"
	"net/http"
)
func current(ctx context.Context) {
	url := "/v1/current"
	Request(ctx, http.MethodPost, url)
}

func legacy(ctx context.Context) {
	Request(ctx, http.MethodPost, "/v1/legacy")
}
func wrongVerb(ctx context.Context) {
	url := "/v1/verb"
	Request(ctx, http.MethodPost, url)
}
`
	if err := os.WriteFile(filepath.Join(sourceDir, "fixture.go"), []byte(source), 0o600); err != nil {
		t.Fatalf("write Go fixture: %v", err)
	}

	specJSON := `{
		"openapi":"3.0.0",
		"paths":{
			"/v1/current":{"post":{"operationId":"Current","deprecated":true}},
			"/v1/spec-only":{"get":{"operationId":"SpecOnly"}},
			"/v1/verb":{"get":{"operationId":"Verb"}}
		}
	}`
	spec, err := LoadOpenAPI(strings.NewReader(specJSON))
	if err != nil {
		t.Fatalf("LoadOpenAPI: %v", err)
	}
	client, err := LoadClientOperations(sourceDir)
	if err != nil {
		t.Fatalf("LoadClientOperations: %v", err)
	}

	report := Compare(client, spec)
	if len(report.Exact) != 1 || report.Exact[0].Path != "/v1/current" {
		t.Fatalf("exact = %+v", report.Exact)
	}
	if len(report.DeprecatedExact) != 1 || report.DeprecatedExact[0].OperationID != "Current" {
		t.Fatalf("deprecated exact = %+v", report.DeprecatedExact)
	}
	if len(report.ClientOnly) != 1 || report.ClientOnly[0].Path != "/v1/legacy" {
		t.Fatalf("client only = %+v", report.ClientOnly)
	}
	if len(report.SpecOnly) != 2 || report.SpecOnly[0].Path != "/v1/spec-only" || report.SpecOnly[1].Path != "/v1/verb" {
		t.Fatalf("spec only = %+v", report.SpecOnly)
	}
	if len(report.MethodMismatches) != 1 || report.MethodMismatches[0].Path != "/v1/verb" || report.MethodMismatches[0].ClientMethod != "POST" || report.MethodMismatches[0].SpecMethods[0] != "GET" {
		t.Fatalf("method mismatches = %+v", report.MethodMismatches)
	}
}

func TestCompareKeepsUnimplementedVerbOnPartiallyImplementedPath(t *testing.T) {
	t.Parallel()

	report := Compare(
		[]Operation{{Method: "GET", Path: "/v1/resource", GoName: "GetResource"}},
		[]Operation{
			{Method: "GET", Path: "/v1/resource", OperationID: "GetResource"},
			{Method: "POST", Path: "/v1/resource", OperationID: "CreateResource"},
		},
	)

	if len(report.Exact) != 1 || len(report.SpecOnly) != 1 {
		t.Fatalf("report = %+v", report)
	}
	if report.SpecOnly[0].Method != "POST" || report.SpecOnly[0].OperationID != "CreateResource" {
		t.Fatalf("spec only = %+v", report.SpecOnly)
	}
}

func TestLoadClientOperationsIgnoresTestFilesAndDeduplicatesMethodPath(t *testing.T) {
	t.Parallel()

	sourceDir := t.TempDir()
	production := `package fixture
import "net/http"
func one() { Request(http.MethodPost, "/v1/same") }
func two() { Request(http.MethodPost, "/v1/same") }
`
	testSource := `package fixture
import "net/http"
func testOnly() { Request(http.MethodDelete, "/v1/test-only") }
`
	if err := os.WriteFile(filepath.Join(sourceDir, "fixture.go"), []byte(production), 0o600); err != nil {
		t.Fatalf("write production fixture: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "fixture_test.go"), []byte(testSource), 0o600); err != nil {
		t.Fatalf("write test fixture: %v", err)
	}

	operations, err := LoadClientOperations(sourceDir)
	if err != nil {
		t.Fatalf("LoadClientOperations: %v", err)
	}
	if len(operations) != 1 || operations[0].Method != "POST" || operations[0].Path != "/v1/same" || operations[0].GoName != "one, two" {
		t.Fatalf("operations = %+v", operations)
	}
}

func TestLoadClientOperationsRecognizesGoDeprecationComments(t *testing.T) {
	t.Parallel()

	sourceDir := t.TempDir()
	source := `package fixture
import "net/http"
// Legacy calls the retired endpoint.
//
// Deprecated: Use Current.
func Legacy() { Request(http.MethodPost, "/v1/legacy") }
func Current() { Request(http.MethodPost, "/v2/current") }
`
	if err := os.WriteFile(filepath.Join(sourceDir, "fixture.go"), []byte(source), 0o600); err != nil {
		t.Fatalf("write Go fixture: %v", err)
	}

	operations, err := LoadClientOperations(sourceDir)
	if err != nil {
		t.Fatalf("LoadClientOperations: %v", err)
	}
	if len(operations) != 2 {
		t.Fatalf("got %d operations, want 2", len(operations))
	}
	if !operations[0].GoDeprecated {
		t.Fatalf("legacy operation should be marked deprecated: %#v", operations[0])
	}
	if operations[1].GoDeprecated {
		t.Fatalf("current operation should not be marked deprecated: %#v", operations[1])
	}
}

func TestLoadClientOperationsResolvesReassignedURLAtEachCall(t *testing.T) {
	t.Parallel()

	sourceDir := t.TempDir()
	source := `package fixture
import "net/http"
func TwoCalls() {
	url := "/v1/one"
	Request(http.MethodPost, url)
	url = "/v1/two"
	Request(http.MethodPost, url)
}
`
	if err := os.WriteFile(filepath.Join(sourceDir, "fixture.go"), []byte(source), 0o600); err != nil {
		t.Fatalf("write Go fixture: %v", err)
	}

	operations, err := LoadClientOperations(sourceDir)
	if err != nil {
		t.Fatalf("LoadClientOperations: %v", err)
	}
	if len(operations) != 2 || operations[0].Path != "/v1/one" || operations[1].Path != "/v1/two" {
		t.Fatalf("operations = %+v", operations)
	}
}

func TestLoadClientOperationsFailsOnUnresolvedRequest(t *testing.T) {
	t.Parallel()

	sourceDir := t.TempDir()
	source := `package fixture
import "net/http"
func Dynamic(path string) { Request(http.MethodPost, path) }
`
	if err := os.WriteFile(filepath.Join(sourceDir, "fixture.go"), []byte(source), 0o600); err != nil {
		t.Fatalf("write Go fixture: %v", err)
	}

	_, err := LoadClientOperations(sourceDir)
	if err == nil || !strings.Contains(err.Error(), "unresolved Request call") {
		t.Fatalf("error = %v, want unresolved Request call", err)
	}
}

func TestWriteMarkdownRendersDeterministicActionableSections(t *testing.T) {
	t.Parallel()

	report := Report{
		Exact:           []Operation{{Method: "POST", Path: "/v1/current", GoName: "Current", File: "current.go"}},
		DeprecatedExact: []Operation{{Method: "POST", Path: "/v1/current", OperationID: "CurrentAPI", Deprecated: true, GoName: "Current", File: "current.go"}},
		ClientOnly:      []Operation{{Method: "POST", Path: "/v1/legacy", GoName: "Legacy", File: "legacy.go"}},
		SpecOnly: []Operation{
			{Method: "GET", Path: "/v1/missing", OperationID: "MissingAPI"},
			{Method: "GET", Path: "/v1/verb", OperationID: "VerbAPI"},
		},
		MethodMismatches: []MethodMismatch{{
			Path: "/v1/verb", ClientMethod: "POST", SpecMethods: []string{"GET"}, GoName: "WrongVerb", File: "verb.go",
		}},
	}
	var output bytes.Buffer
	if err := WriteMarkdown(&output, report, Metadata{
		SwaggerFile:   "swagger.json",
		SwaggerSHA256: "abc123",
		GeneratedDate: "2026-09-02",
	}); err != nil {
		t.Fatalf("WriteMarkdown: %v", err)
	}

	want := `# Ozon Seller API contract audit

- Generated: 2026-09-02
- Swagger: ` + "`swagger.json`" + `
- Swagger SHA-256: ` + "`abc123`" + `
- Client operations: 3
- Swagger operations: 3
- Exact method/path matches: 1
- Client-only paths: 1
- Method mismatches: 1
- Swagger-only operations: 2
- Deprecated exact matches: 1
- Deprecated matches missing Go doc: 1

## Client paths absent from Swagger

| Method | Path | Go method | Source |
|---|---|---|---|
| POST | ` + "`/v1/legacy`" + ` | ` + "`Legacy`" + ` | ` + "`legacy.go`" + ` |

## Method mismatches

| Path | Client | Swagger | Go method | Source |
|---|---|---|---|---|
| ` + "`/v1/verb`" + ` | POST | GET | ` + "`WrongVerb`" + ` | ` + "`verb.go`" + ` |

## Deprecated exact matches

| Method | Path | Go method | Go deprecated | Operation ID | Source |
|---|---|---|---|---|---|
| POST | ` + "`/v1/current`" + ` | ` + "`Current`" + ` | no | ` + "`CurrentAPI`" + ` | ` + "`current.go`" + ` |

## Swagger operations not implemented by the client

| Method | Path | Operation ID |
|---|---|---|
| GET | ` + "`/v1/missing`" + ` | ` + "`MissingAPI`" + ` |
| GET | ` + "`/v1/verb`" + ` | ` + "`VerbAPI`" + ` |
`
	if output.String() != want {
		t.Fatalf("markdown output:\n%s\nwant:\n%s", output.String(), want)
	}
}
