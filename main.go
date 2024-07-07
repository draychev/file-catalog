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

	"github.com/spf13/cobra"
)

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

func main() {
	var storagePath string
	var outputPath string

	var rootCmd = &cobra.Command{
		Use:   "filehasher",
		Short: "File Hasher is a tool to hash files and store their metadata.",
	}

	var hashCmd = &cobra.Command{
		Use:   "hash",
		Short: "Hash files in the specified directory.",
		Run: func(cmd *cobra.Command, args []string) {
			files, err := ioutil.ReadDir(storagePath)
			if err != nil {
				fmt.Println("Error reading directory:", err)
				return
			}

			var metas []FileMeta
			for _, file := range files {
				if !file.IsDir() {
					meta, err := getFileMeta(filepath.Join(storagePath, file.Name()))
					if err != nil {
						fmt.Println("Error getting file metadata:", err)
						continue
					}
					metas = append(metas, meta)
				}
			}

			if err := serializeFileMeta(metas, outputPath); err != nil {
				fmt.Println("Error serializing metadata:", err)
			} else {
				fmt.Println("File metadata has been serialized to", outputPath)
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
				fmt.Println("Error deserializing metadata:", err)
				return
			}

			for _, meta := range metas {
				fmt.Printf("File: %s, Hash: %s, Size: %d bytes, Created At: %s, Created By: %s, Last Modified: %s, Modified By: %s, Last Accessed: %s, Accessed By: %s\n",
					meta.FileName, meta.Hash, meta.FileSize, meta.CreatedAt, meta.CreatedBy, meta.LastModified, meta.ModifiedBy, meta.LastAccessed, meta.AccessedBy)
			}
		},
	}

	showCmd.Flags().StringVarP(&outputPath, "input", "i", "file_metadata.json", "Path to the input file")

	rootCmd.AddCommand(hashCmd, showCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error executing command:", err)
	}
}
