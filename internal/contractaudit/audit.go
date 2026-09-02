package contractaudit

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Operation struct {
	Method       string
	Path         string
	OperationID  string
	Deprecated   bool
	GoDeprecated bool
	GoName       string
	File         string
}

type MethodMismatch struct {
	Path         string
	ClientMethod string
	SpecMethods  []string
	GoName       string
	File         string
}

type Report struct {
	Exact            []Operation
	DeprecatedExact  []Operation
	ClientOnly       []Operation
	SpecOnly         []Operation
	MethodMismatches []MethodMismatch
}

type Metadata struct {
	SwaggerFile   string
	SwaggerSHA256 string
	GeneratedDate string
}

func WriteMarkdown(writer io.Writer, report Report, metadata Metadata) error {
	clientOperations := len(report.Exact) + len(report.ClientOnly) + len(report.MethodMismatches)
	swaggerOperations := len(report.Exact) + len(report.SpecOnly)
	undocumentedDeprecated := 0
	for _, operation := range report.DeprecatedExact {
		if !operation.GoDeprecated {
			undocumentedDeprecated++
		}
	}
	for _, mismatch := range report.MethodMismatches {
		swaggerOperations += len(mismatch.SpecMethods)
	}

	if _, err := fmt.Fprintf(writer, `# Ozon Seller API contract audit

- Generated: %s
- Swagger: `+"`%s`"+`
- Swagger SHA-256: `+"`%s`"+`
- Client operations: %d
- Swagger operations: %d
- Exact method/path matches: %d
- Client-only paths: %d
- Method mismatches: %d
- Swagger-only operations: %d
- Deprecated exact matches: %d
- Deprecated matches missing Go doc: %d

## Client paths absent from Swagger

| Method | Path | Go method | Source |
|---|---|---|---|
`, metadata.GeneratedDate, metadata.SwaggerFile, metadata.SwaggerSHA256, clientOperations, swaggerOperations, len(report.Exact), len(report.ClientOnly), len(report.MethodMismatches), len(report.SpecOnly), len(report.DeprecatedExact), undocumentedDeprecated); err != nil {
		return err
	}
	for _, operation := range report.ClientOnly {
		if _, err := fmt.Fprintf(writer, "| %s | `%s` | `%s` | `%s` |\n", operation.Method, operation.Path, operation.GoName, operation.File); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprint(writer, `
## Method mismatches

| Path | Client | Swagger | Go method | Source |
|---|---|---|---|---|
`); err != nil {
		return err
	}
	for _, mismatch := range report.MethodMismatches {
		if _, err := fmt.Fprintf(writer, "| `%s` | %s | %s | `%s` | `%s` |\n", mismatch.Path, mismatch.ClientMethod, strings.Join(mismatch.SpecMethods, ", "), mismatch.GoName, mismatch.File); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprint(writer, `
## Deprecated exact matches

| Method | Path | Go method | Go deprecated | Operation ID | Source |
|---|---|---|---|---|---|
`); err != nil {
		return err
	}
	for _, operation := range report.DeprecatedExact {
		goDeprecated := "no"
		if operation.GoDeprecated {
			goDeprecated = "yes"
		}
		if _, err := fmt.Fprintf(writer, "| %s | `%s` | `%s` | %s | `%s` | `%s` |\n", operation.Method, operation.Path, operation.GoName, goDeprecated, operation.OperationID, operation.File); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprint(writer, `
## Swagger operations not implemented by the client

| Method | Path | Operation ID |
|---|---|---|
`); err != nil {
		return err
	}
	for _, operation := range report.SpecOnly {
		if _, err := fmt.Fprintf(writer, "| %s | `%s` | `%s` |\n", operation.Method, operation.Path, operation.OperationID); err != nil {
			return err
		}
	}
	return nil
}

func LoadOpenAPI(reader io.Reader) ([]Operation, error) {
	var document struct {
		Paths map[string]map[string]struct {
			OperationID string `json:"operationId"`
			Deprecated  bool   `json:"deprecated"`
		} `json:"paths"`
	}
	if err := json.NewDecoder(reader).Decode(&document); err != nil {
		return nil, fmt.Errorf("decode OpenAPI document: %w", err)
	}

	var operations []Operation
	for path, pathItem := range document.Paths {
		for method, operation := range pathItem {
			method = strings.ToUpper(method)
			if !isHTTPMethod(method) {
				continue
			}
			operations = append(operations, Operation{
				Method:      method,
				Path:        path,
				OperationID: operation.OperationID,
				Deprecated:  operation.Deprecated,
			})
		}
	}
	sortOperations(operations)
	return operations, nil
}

func LoadClientOperations(root string) ([]Operation, error) {
	operationByKey := make(map[string]Operation)
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
			return nil
		}
		fileSet := token.NewFileSet()
		file, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		relativePath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		for _, declaration := range file.Decls {
			function, ok := declaration.(*ast.FuncDecl)
			if !ok || function.Body == nil {
				continue
			}
			stringValues := collectStringValues(function.Body)
			goDeprecated := function.Doc != nil && strings.Contains(function.Doc.Text(), "Deprecated:")
			ast.Inspect(function.Body, func(node ast.Node) bool {
				call, ok := node.(*ast.CallExpr)
				if !ok || !isRequestCall(call.Fun) {
					return true
				}
				method, path := requestMethodAndPath(call.Args, stringValues)
				if method == "" || path == "" {
					return true
				}
				key := method + " " + path
				if existing, exists := operationByKey[key]; !exists {
					operationByKey[key] = Operation{
						Method:       method,
						Path:         path,
						GoDeprecated: goDeprecated,
						GoName:       function.Name.Name,
						File:         filepath.ToSlash(relativePath),
					}
				} else {
					existing.GoName = mergeSortedValues(existing.GoName, function.Name.Name)
					existing.File = mergeSortedValues(existing.File, filepath.ToSlash(relativePath))
					existing.GoDeprecated = existing.GoDeprecated && goDeprecated
					operationByKey[key] = existing
				}
				return true
			})
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk Go client sources: %w", err)
	}

	operations := make([]Operation, 0, len(operationByKey))
	for _, operation := range operationByKey {
		operations = append(operations, operation)
	}
	sortOperations(operations)
	return operations, nil
}

func mergeSortedValues(existing, value string) string {
	values := strings.Split(existing, ", ")
	for _, candidate := range values {
		if candidate == value {
			return existing
		}
	}
	values = append(values, value)
	sort.Strings(values)
	return strings.Join(values, ", ")
}

func Compare(client, spec []Operation) Report {
	specByKey := make(map[string]Operation, len(spec))
	specMethodsByPath := make(map[string][]string)
	for _, operation := range spec {
		specByKey[operation.Method+" "+operation.Path] = operation
		specMethodsByPath[operation.Path] = append(specMethodsByPath[operation.Path], operation.Method)
	}
	for path := range specMethodsByPath {
		sort.Strings(specMethodsByPath[path])
	}

	clientPaths := make(map[string]struct{})
	matchedSpec := make(map[string]struct{})
	var report Report
	for _, operation := range client {
		clientPaths[operation.Path] = struct{}{}
		key := operation.Method + " " + operation.Path
		if specOperation, exists := specByKey[key]; exists {
			operation.OperationID = specOperation.OperationID
			operation.Deprecated = specOperation.Deprecated
			report.Exact = append(report.Exact, operation)
			if operation.Deprecated {
				report.DeprecatedExact = append(report.DeprecatedExact, operation)
			}
			matchedSpec[key] = struct{}{}
			continue
		}
		if methods, pathExists := specMethodsByPath[operation.Path]; pathExists {
			report.MethodMismatches = append(report.MethodMismatches, MethodMismatch{
				Path:         operation.Path,
				ClientMethod: operation.Method,
				SpecMethods:  append([]string(nil), methods...),
				GoName:       operation.GoName,
				File:         operation.File,
			})
			continue
		}
		report.ClientOnly = append(report.ClientOnly, operation)
	}
	for _, operation := range spec {
		if _, matched := matchedSpec[operation.Method+" "+operation.Path]; matched {
			continue
		}
		if _, clientPathExists := clientPaths[operation.Path]; clientPathExists {
			continue
		}
		report.SpecOnly = append(report.SpecOnly, operation)
	}

	sortOperations(report.Exact)
	sortOperations(report.DeprecatedExact)
	sortOperations(report.ClientOnly)
	sortOperations(report.SpecOnly)
	sort.Slice(report.MethodMismatches, func(i, j int) bool {
		if report.MethodMismatches[i].Path != report.MethodMismatches[j].Path {
			return report.MethodMismatches[i].Path < report.MethodMismatches[j].Path
		}
		return report.MethodMismatches[i].ClientMethod < report.MethodMismatches[j].ClientMethod
	})
	return report
}

func collectStringValues(body *ast.BlockStmt) map[string]string {
	values := make(map[string]string)
	ast.Inspect(body, func(node ast.Node) bool {
		switch statement := node.(type) {
		case *ast.AssignStmt:
			for i, left := range statement.Lhs {
				if i >= len(statement.Rhs) {
					continue
				}
				identifier, ok := left.(*ast.Ident)
				if !ok {
					continue
				}
				if value := stringLiteral(statement.Rhs[i], nil); value != "" {
					values[identifier.Name] = value
				}
			}
		case *ast.ValueSpec:
			for i, name := range statement.Names {
				if i >= len(statement.Values) {
					continue
				}
				if value := stringLiteral(statement.Values[i], nil); value != "" {
					values[name.Name] = value
				}
			}
		}
		return true
	})
	return values
}

func isRequestCall(expression ast.Expr) bool {
	switch function := expression.(type) {
	case *ast.Ident:
		return function.Name == "Request"
	case *ast.SelectorExpr:
		return function.Sel.Name == "Request"
	default:
		return false
	}
}

func requestMethodAndPath(arguments []ast.Expr, stringsByName map[string]string) (string, string) {
	var method, path string
	for _, argument := range arguments {
		if method == "" {
			method = httpMethod(argument)
		}
		if path == "" {
			candidate := stringLiteral(argument, stringsByName)
			if strings.HasPrefix(candidate, "/") {
				path = candidate
			}
		}
	}
	return method, path
}

func httpMethod(expression ast.Expr) string {
	switch value := expression.(type) {
	case *ast.SelectorExpr:
		if identifier, ok := value.X.(*ast.Ident); ok && identifier.Name == "http" && strings.HasPrefix(value.Sel.Name, "Method") {
			return strings.ToUpper(strings.TrimPrefix(value.Sel.Name, "Method"))
		}
	case *ast.BasicLit:
		if value.Kind == token.STRING {
			literal, err := strconv.Unquote(value.Value)
			if err == nil && isHTTPMethod(strings.ToUpper(literal)) {
				return strings.ToUpper(literal)
			}
		}
	}
	return ""
}

func stringLiteral(expression ast.Expr, stringsByName map[string]string) string {
	switch value := expression.(type) {
	case *ast.BasicLit:
		if value.Kind != token.STRING {
			return ""
		}
		literal, err := strconv.Unquote(value.Value)
		if err == nil {
			return literal
		}
	case *ast.Ident:
		return stringsByName[value.Name]
	}
	return ""
}

func isHTTPMethod(method string) bool {
	switch method {
	case "GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS", "TRACE":
		return true
	default:
		return false
	}
}

func sortOperations(operations []Operation) {
	sort.Slice(operations, func(i, j int) bool {
		if operations[i].Path != operations[j].Path {
			return operations[i].Path < operations[j].Path
		}
		if operations[i].Method != operations[j].Method {
			return operations[i].Method < operations[j].Method
		}
		return operations[i].GoName < operations[j].GoName
	})
}
