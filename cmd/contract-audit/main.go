package main

import (
	"bytes"
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/ucoms-dev/ozon-api-client/internal/contractaudit"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, output io.Writer) error {
	flags := flag.NewFlagSet("contract-audit", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	var swaggerPath string
	var clientPath string
	var generatedDate string
	var outputPath string
	flags.StringVar(&swaggerPath, "swagger", "", "path to the Ozon OpenAPI JSON document")
	flags.StringVar(&clientPath, "client", "ozon", "path to the Go client sources")
	flags.StringVar(&generatedDate, "date", time.Now().UTC().Format(time.DateOnly), "audit date in YYYY-MM-DD format")
	flags.StringVar(&outputPath, "output", "", "write the Markdown report to this file instead of stdout")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if swaggerPath == "" {
		return fmt.Errorf("-swagger is required")
	}

	swaggerData, err := os.ReadFile(swaggerPath)
	if err != nil {
		return fmt.Errorf("read Swagger document: %w", err)
	}
	spec, err := contractaudit.LoadOpenAPI(bytes.NewReader(swaggerData))
	if err != nil {
		return err
	}
	client, err := contractaudit.LoadClientOperations(clientPath)
	if err != nil {
		return err
	}
	digest := sha256.Sum256(swaggerData)
	metadata := contractaudit.Metadata{
		SwaggerFile:   filepath.Base(swaggerPath),
		SwaggerSHA256: fmt.Sprintf("%x", digest),
		GeneratedDate: generatedDate,
	}
	if outputPath == "" {
		return contractaudit.WriteMarkdown(output, contractaudit.Compare(client, spec), metadata)
	}

	var report bytes.Buffer
	if err := contractaudit.WriteMarkdown(&report, contractaudit.Compare(client, spec), metadata); err != nil {
		return err
	}
	if err := os.WriteFile(outputPath, report.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write audit report: %w", err)
	}
	return nil
}
