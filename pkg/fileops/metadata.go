package fileops

import (
	"os"
	"time"
)

// GetMetadata gathers metadata for a fileops.
func GetMetadata(filePath string) (FileMeta, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return FileMeta{}, err
	}

	hash, err := HashSHA256(filePath)
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
