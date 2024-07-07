package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/draychev/file-catalog/pkg/fileops"
	"github.com/draychev/file-catalog/pkg/storage"
)

// GetDupes returns the Cobra Command to find duplicated files
func GetDupes(outputPath string) *cobra.Command {
	var dupesCmd = &cobra.Command{
		Use:   "dupes",
		Short: "Find and display fileops with identical hashes.",
		Run: func(cmd *cobra.Command, args []string) {
			metas, err := storage.DeserializeFileMeta(outputPath)
			if err != nil {
				log.Error().Err(err).Msgf("Error deserializing metadata from: %s", outputPath)
				return
			}

			hashMap := make(map[string][]fileops.FileMeta)

			for _, meta := range metas {
				hashMap[meta.Hash] = append(hashMap[meta.Hash], meta)
			}

			fmt.Printf("%-64s %-30s %-25s %-30s %-25s\n", "Hash", "First File", "Created At", "Second File", "Created At")
			for hash, files := range hashMap {
				if len(files) > 1 {
					for i := 0; i < len(files)-1; i++ {
						for j := i + 1; j < len(files); j++ {
							fmt.Printf("%-64s %-30s %-25s %-30s %-25s\n",
								hash, files[i].FileName, files[i].CreatedAt, files[j].FileName, files[j].CreatedAt)
						}
					}
				}
			}
		},
	}

	dupesCmd.Flags().StringVarP(&outputPath, "input", "i", "file_metadata.json", "Path to the input fileops")

	return dupesCmd
}
