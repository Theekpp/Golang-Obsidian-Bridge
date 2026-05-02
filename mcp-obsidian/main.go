package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/user/mcp-obsidian/internal/obsidian"
	"github.com/user/mcp-obsidian/internal/tools"
)

const (
	defaultBaseURL = "http://127.0.0.1:27123"
	version        = "1.0.0"
)

func main() {
	var (
		baseURL = flag.String("url", envOrDefault("OBSIDIAN_URL", defaultBaseURL),
			"Obsidian Local REST API base URL (env: OBSIDIAN_URL)")
		apiKey = flag.String("key", os.Getenv("OBSIDIAN_API_KEY"),
			"Obsidian Local REST API key (env: OBSIDIAN_API_KEY)")
		showVersion = flag.Bool("version", false, "Print version and exit")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("mcp-obsidian v%s\n", version)
		os.Exit(0)
	}

	if *apiKey == "" {
		log.Fatal("ERROR: Obsidian API key is required.\n" +
			"Set the OBSIDIAN_API_KEY environment variable or use the -key flag.\n" +
			"You can find your API key in Obsidian → Settings → Local REST API.")
	}

	obsidianClient := obsidian.NewClient(*baseURL, *apiKey)

	mcpServer := server.NewMCPServer(
		"mcp-obsidian",
		version,
		server.WithToolCapabilities(true),
	)

	tools.Register(mcpServer, obsidianClient)

	log.Printf("mcp-obsidian v%s starting (Obsidian at %s)", version, *baseURL)

	if err := server.ServeStdio(mcpServer); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
