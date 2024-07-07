package commands

import (
	"fmt"
	"runtime"
	"sync"

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

			numCPU := runtime.NumCPU()
			var wg sync.WaitGroup
			fileChan := make(chan string, len(files))
			metaChan := make(chan fileops.FileMeta, len(files))

			wg.Add(numCPU)

			for i := 0; i < numCPU; i++ {
				go func() {
					defer wg.Done()
					for file := range fileChan {
						meta, err := fileops.GetMetadata(file)
						if err != nil {
							log.Error().Err(err).Msgf("Error getting fileops metadata: %s", err)
							continue
						}
						metaChan <- meta
					}
				}()
			}

			go func() {
				for _, file := range files {
					fileChan <- file
				}
				close(fileChan)
				wg.Wait()
				close(metaChan)
			}()

			var metas []fileops.FileMeta
			totalFiles := len(files)
			processedFiles := 0

			for meta := range metaChan {
				metas = append(metas, meta)
				processedFiles++
				percentComplete := float64(processedFiles) / float64(totalFiles) * 100
				fmt.Printf("%.2f%% - Hashing fileops: %s\n", percentComplete, meta.FileName)
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
