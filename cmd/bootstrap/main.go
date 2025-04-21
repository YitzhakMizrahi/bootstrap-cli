package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/core/config"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/core/system"
	"github.com/YitzhakMizrahi/bootstrap-cli/pkg/plugin/manager"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	// Initialize configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize system information
	sysInfo, err := system.GetSystemInfo()
	if err != nil {
		fmt.Printf("Error getting system information: %v\n", err)
		os.Exit(1)
	}

	// Initialize plugin manager
	pluginManager := manager.New()

	// TODO: Implement main bootstrap logic
	fmt.Printf("Bootstrap CLI initialized for %s\n", sysInfo.OS)
} 