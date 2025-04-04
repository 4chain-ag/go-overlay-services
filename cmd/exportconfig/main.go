package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"slices"
	"strings"

	config "github.com/4chain-ag/go-overlay-services/pkg/appconfig"
	"github.com/google/uuid"
)

func main() {
	regenToken := flag.Bool("regen-token", false, "Regenerate admin bearer token")
	flag.BoolVar(regenToken, "t", false, "Regenerate admin bearer token (shorthand)")

	outputFile := flag.String("output-file", "config.yaml", "Output configuration file path")
	flag.StringVar(outputFile, "o", "config.yaml", "Output configuration file path (shorthand)")

	flag.Parse()

	cfg := config.Defaults()

	if *regenToken {
		cfg.AdminBearerToken = uuid.NewString()
	}

	ext := strings.TrimPrefix(filepath.Ext(*outputFile), ".")

	if !isSupportedExt(ext) {
		log.Fatalf("Unsupported output file extension: %s", ext)
	}

	var err error
	switch ext {
	case "json":
		err = config.ToJSONFile(&cfg, *outputFile)
	case "env", "dotenv":
		err = config.ToEnvFile(&cfg, *outputFile)
	default: // yaml, yml
		err = config.ToYAMLFile(&cfg, *outputFile)
	}

	if err != nil {
		log.Fatalf("Error writing configuration: %v\n", err)
	}

	fmt.Printf("Configuration written to %s\n", *outputFile)
}

func isSupportedExt(ext string) bool {
	return slices.Contains(config.SupportedExts(), ext)
}
