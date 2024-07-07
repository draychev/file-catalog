package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/draychev/go-toolbox/pkg/logger"
	"github.com/spf13/cobra"
)

var log = logger.NewPretty("file-hasher")

// FileMeta holds metadata of a file.
type FileMeta struct {
	FileName     string `json:"file_name"`
	Hash         string `json:"hash"`
	FileSize     int64  `json:"file_size"`
	CreatedAt    string `json:"created_at"`
	CreatedBy    string `json:"created_by"`
	LastModified string `json:"last_modified"`
	ModifiedBy   string `json:"modified_by"`
	LastAccessed string `json:"last_accessed"`
	AccessedBy   string `json:"accessed_by"`
}

// hashFile computes the SHA256 hash of a file.
func hashFile(filePath string) (string, error) {
	var hasher hash.Hash = sha256.New()
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

// getFileMeta gathers metadata for a file.
func getFileMeta(filePath string) (FileMeta, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return FileMeta{}, err
	}

	hash, err := hashFile(filePath)
	if err != nil {
		return FileMeta{}, err
	}

	// Replace with actual owner info retrieval logic if needed.
	createdBy := "unknown"
	modifiedBy := "unknown"
	accessedBy := "unknown"

	return FileMeta{
		FileName:     filePath,
		Hash:         hash,
		FileSize:     info.Size(),
		CreatedAt:    info.ModTime().Format(time.RFC3339),
		CreatedBy:    createdBy,
		LastModified: info.ModTime().Format(time.RFC3339),
		ModifiedBy:   modifiedBy,
		LastAccessed: info.ModTime().Format(time.RFC3339), // Placeholder
		AccessedBy:   accessedBy,                          // Placeholder
	}, nil
}

// serializeFileMeta serializes the file metadata to a JSON file.
func serializeFileMeta(metas []FileMeta, outputPath string) error {
	data, err := json.MarshalIndent(metas, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(outputPath, data, 0644)
}

// deserializeFileMeta deserializes the file metadata from a JSON file.
func deserializeFileMeta(inputPath string) ([]FileMeta, error) {
	data, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return nil, err
	}

	var metas []FileMeta
	if err := json.Unmarshal(data, &metas); err != nil {
		return nil, err
	}

	return metas, nil
}

// collectFiles recursively collects files in the specified directory.
func collectFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func main() {
	var storagePath string
	var outputPath string

	var rootCmd = &cobra.Command{
		Use:   "file-catalog",
		Short: "File Catalog is a tool to hash files and store their metadata.",
	}

	var hashCmd = &cobra.Command{
		Use:   "hash",
		Short: "Hash files in the specified directory.",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info().Msgf("Hashing files in directory %s", storagePath)
			files, err := collectFiles(storagePath)
			if err != nil {
				log.Error().Err(err).Msgf("Error reading directory: %s", err)
				return
			}

			log.Info().Msgf("Here are the files we found in %s: %+v", storagePath, files)

			var metas []FileMeta
			totalFiles := len(files)

			for i, file := range files {
				meta, err := getFileMeta(file)
				if err != nil {
					log.Error().Err(err).Msgf("Error getting file metadata: %s", err)
					continue
				}
				metas = append(metas, meta)
				percentComplete := float64(i+1) / float64(totalFiles) * 100
				fmt.Printf("%.2f%% - Hashing file: %s\n", percentComplete, file)
			}

			if err := serializeFileMeta(metas, outputPath); err != nil {
				log.Error().Err(err).Msgf("Error serializing metadata: %s", err)
			} else {
				log.Info().Msgf("File metadata has been serialized to %s", outputPath)
			}
		},
	}

	hashCmd.Flags().StringVarP(&storagePath, "storage", "s", "/storage", "Path to the storage directory")
	hashCmd.Flags().StringVarP(&outputPath, "output", "o", "file_metadata.json", "Path to the output file")

	var showCmd = &cobra.Command{
		Use:   "show",
		Short: "Show hashed files from the serialized metadata file.",
		Run: func(cmd *cobra.Command, args []string) {
			metas, err := deserializeFileMeta(outputPath)
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

	showCmd.Flags().StringVarP(&outputPath, "input", "i", "file_metadata.json", "Path to the input file")

	var dupesCmd = &cobra.Command{
		Use:   "dupes",
		Short: "Find and display files with identical hashes.",
		Run: func(cmd *cobra.Command, args []string) {
			metas, err := deserializeFileMeta(outputPath)
			if err != nil {
				log.Error().Err(err).Msgf("Error deserializing metadata from: %s", outputPath)
				return
			}

			hashMap := make(map[string][]FileMeta)

			for _, meta := range metas {
				hashMap[meta.Hash] = append(hashMap[meta.Hash], meta)
			}

			fmt.Printf("%-64s %-30s %-25s %-30s %-25s\n", "Hash", "First File Name", "First Created At", "Second File Name", "Second Created At")
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

	dupesCmd.Flags().StringVarP(&outputPath, "input", "i", "file_metadata.json", "Path to the input file")

	rootCmd.AddCommand(hashCmd, showCmd, dupesCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err).Msgf("Error executing command: %s", err)
	}
}
