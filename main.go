// Copyright (c) 2025 Tenebris Technologies Inc.
// This software is licensed under the MIT License (see LICENSE for details).

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/PivotLLM/MCPLaunchPad/example1"
	"github.com/PivotLLM/MCPLaunchPad/example2"
	"github.com/PivotLLM/MCPLaunchPad/global"
	"github.com/PivotLLM/MCPLaunchPad/mcpserver"
	"github.com/PivotLLM/MCPLaunchPad/mlogger"
)

// Version information
const (
	AppName    = "MCPLaunchPad"
	AppVersion = "0.0.2"
)

func main() {
	var err error
	var listen string

	// Define command line flags
	debugFlag := flag.Bool("debug", true, "Enable debug mode")
	portFlag := flag.Int("port", 8888, "Port to listen on")
	helpFlag := flag.Bool("help", false, "Show help information")
	versionFlag := flag.Bool("version", false, "Show version information")

	// Set custom usage message
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("  %s [options]\n\n", os.Args[0])
		fmt.Printf("Options:\n")
		flag.PrintDefaults()
	}

	// Parse command line flags
	flag.Parse()

	// Show help and exit if requested
	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	// Show version and exit if requested
	if *versionFlag {
		fmt.Printf("%s version %s\n", AppName, AppVersion)
		os.Exit(0)
	}

	// Use the flag values
	debug := *debugFlag
	if *portFlag > 0 && *portFlag < 65536 {
		listen = fmt.Sprintf("localhost:%d", *portFlag)
	} else {
		listen = "localhost:8888"
	}

	logger, err := mlogger.New(
		mlogger.WithPrefix("MCP"),
		mlogger.WithDateFormat("2006-01-02 15:04:05"),
		mlogger.WithLogFile("mcp.log"),
		mlogger.WithLogStdout(true),
		mlogger.WithDebug(debug),
	)
	if err != nil {
		fmt.Printf("Unable to create logger: %v", err)
		os.Exit(1)
	}

	// Get the user's home directory if possible
	homeDir, err := os.UserHomeDir()
	if err == nil {
		envFile := homeDir + string(os.PathSeparator) + ".mcp"
		err = godotenv.Load(envFile)
		if err == nil {
			logger.Infof("Loaded environment variables from %s", envFile)
		}
	}

	// Load BaseURL and auth key from environment variables`
	APIBaseURL := os.Getenv("API_BASE_URL")
	if APIBaseURL == "" {
		//logger.Fatalf("API_BASE_URL environment variable is not set")
		//os.Exit(1)
		APIBaseURL = "https://api.example.com"
	}
	APIAuthKey := os.Getenv("API_AUTH_KEY")
	if APIAuthKey == "" {
		//logger.Fatalf("API_AUTH_KEY environment variable is not set")
		//os.Exit(1)
		APIAuthKey = "1234567890"
	}
	APIAuthHeader := os.Getenv("API_AUTH_HEADER")
	if APIAuthHeader == "" {
		//logger.Fatalf("API_AUTH_HEADER environment variable is not set")
		//os.Exit(1)
		APIAuthHeader = "X-API-Key"
	}

	// Create the example1 API tool provider with a hard-coded base URL
	tp1 := example1.New(
		example1.WithBaseURL(APIBaseURL),
		example1.WithLogger(logger),
		example1.WithAuthHeader(APIAuthHeader),
		example1.WithAuthKey(APIAuthKey),
	)

	// Create the example2 time tool provider
	tp2 := example2.New(
		example2.WithLogger(logger),
	)

	// Create a slice (list) of tool providers
	providers := []global.ToolProvider{
		tp1,
		tp2,
	}

	// Create MCP server
	mcp, err := mcpserver.New(
		mcpserver.WithListen(listen),
		mcpserver.WithDebug(debug),
		mcpserver.WithLogger(logger),
		mcpserver.WithName(AppName),
		mcpserver.WithVersion(AppVersion),
		mcpserver.WithToolProviders(providers),
	)
	if err != nil {
		logger.Fatalf("Unable to create MCP server: %v", err)
		os.Exit(1)
	}

	// Start MCP server
	if err = mcp.Start(); err != nil {
		logger.Fatalf("MCP server failed to start: %v", err)
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal
	<-sigChan
	logger.Infof("Shutting down...")

	// Stop the MCP server
	if err = mcp.Stop(); err != nil {
		logger.Errorf("Error stopping MCP server: %s", err.Error())
		os.Exit(1)
	}

	logger.Infof("MCP server stopped successfully")

	// Exit with success
	os.Exit(0)
}
