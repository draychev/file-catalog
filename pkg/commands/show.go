package commands

import (
	"fmt"
	"github.com/draychev/file-catalog/pkg/storage"
	"github.com/spf13/cobra"
)

func GetShow(outputPath string) *cobra.Command {
	var showCmd = &cobra.Command{
		Use:   "show",
		Short: "Show hashed fileops from the serialized metadata fileops.",
		Run: func(cmd *cobra.Command, args []string) {
			metas, err := storage.DeserializeFileMeta(outputPath)
			if err != nil {
				log.Error().Err(err).Msgf("Error deserializing metadata from: %s", outputPath)
				return
			}

			for _, meta := range metas {
				fmt.Printf("File: %s, Hash: %s, Size: %d bytes, Created At: %s, Created By: %s, Last Modified: %s, Modified By: %s, Last Accessed: %s, Accessed By: %s\n",
					meta.FileName, meta.Hash, meta.FileSize, meta.CreatedAt, meta.CreatedBy, meta.LastModified, meta.ModifiedBy, meta.LastAccessed, meta.AccessedBy)
			}
		},
	}

	showCmd.Flags().StringVarP(&outputPath, "input", "i", "file_metadata.json", "Path to the input fileops")

	return showCmd
}
