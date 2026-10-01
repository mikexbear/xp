package main

import "flag"

func main() {}

func parseArgs() string {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "Path to YAML configuration")

	flag.Parse()

	return configPath
}
