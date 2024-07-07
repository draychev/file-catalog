package fileops

import (
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"os"
)

// Hash computes the SHA256 hash of a fileops.
func Hash(filePath string) (string, error) {
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
