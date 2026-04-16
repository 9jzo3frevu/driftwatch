package main

import (
	"fmt"
	"os"

	"github.com/driftwatch/internal/config"
	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/report"
	"github.com/driftwatch/internal/source"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfgPath := "driftwatch.yaml"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	loader := source.NewLoader()
	declared, err := loader.LoadFile(cfg.DeclaredFile)
	if err != nil {
		return fmt.Errorf("loading declared state from %q: %w", cfg.DeclaredFile, err)
	}

	fetcher := source.NewHTTPFetcher(cfg.Timeout)
	live, err := fetcher.Fetch(cfg.LiveURL)
	if err != nil {
		return fmt.Errorf("fetching live state from %q: %w", cfg.LiveURL, err)
	}

	results := drift.Detect(declared, live)

	writer, err := report.NewFileWriter(cfg.OutputFile)
	if err != nil {
		return fmt.Errorf("creating report writer: %w", err)
	}

	w, err := writer.Open()
	if err != nil {
		return fmt.Errorf("opening report output: %w", err)
	}
	defer w.Close()

	fmt.Fprintf(os.Stderr, "detected %d drift(s)\n", len(results))

	fmt.Fprintln(os.Stderr, drift.DriftResult(results).Summary())

	formatter := report.NewFormatter(w, cfg.Format)
	if err := formatter.Write(results); err != nil {
		return fmt.Errorf("writing report: %w", err)
	}

	if len(results) > 0 {
		os.Exit(2)
	}
	return nil
}
