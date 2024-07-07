package main

import (
	"github.com/spf13/cobra"

	"github.com/draychev/file-catalog/pkg/commands"
	"github.com/draychev/go-toolbox/pkg/logger"
)

var log = logger.NewPretty("fileops-hasher")

func main() {
	var storagePath string
	var outputPath string

	var rootCmd = &cobra.Command{
		Use:   "fileops-catalog",
		Short: "File Catalog is a tool to hash fileops and store their metadata.",
	}

	rootCmd.AddCommand(
		commands.GetHash(outputPath, storagePath),
		commands.GetShow(outputPath),
		commands.GetDupes(outputPath))
	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err).Msgf("Error executing command: %s", err)
	}
}
