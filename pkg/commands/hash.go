package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/draychev/file-catalog/pkg/fileops"
	"github.com/draychev/file-catalog/pkg/storage"
)

func GetHash(outputPath, storagePath string) *cobra.Command {
	var hashCmd = &cobra.Command{
		Use:   "hash",
		Short: "Hash fileops in the specified directory.",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info().Msgf("Hashing fileops in directory %s", storagePath)
			files, err := fileops.CollectFiles(storagePath)
			if err != nil {
				log.Error().Err(err).Msgf("Error reading directory: %s", err)
				return
			}

			log.Info().Msgf("Here are the fileops we found in %s: %+v", storagePath, files)

			var metas []fileops.FileMeta
			totalFiles := len(files)

			for i, file := range files {
				meta, err := fileops.GetMetadata(file)
				if err != nil {
					log.Error().Err(err).Msgf("Error getting fileops metadata: %s", err)
					continue
				}
				metas = append(metas, meta)
				percentComplete := float64(i+1) / float64(totalFiles) * 100
				fmt.Printf("%.2f%% - Hashing fileops: %s\n", percentComplete, file)
			}

			if err := storage.SerializeFileMeta(metas, outputPath); err != nil {
				log.Error().Err(err).Msgf("Error serializing metadata: %s", err)
			} else {
				log.Info().Msgf("File metadata has been serialized to %s", outputPath)
			}
		},
	}

	hashCmd.Flags().StringVarP(&storagePath, "storage", "s", "/storage", "Path to the storage directory")
	hashCmd.Flags().StringVarP(&outputPath, "output", "o", "file_metadata.json", "Path to the output fileops")

	return hashCmd
}
