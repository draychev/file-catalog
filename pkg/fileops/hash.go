package fileops

import (
	"crypto/sha256"
	"fmt"
	"hash/crc32"
	"io"
	"os"
)

// HashCRC32 computes the CRC32 hash of a fileops.
func HashCRC32(filePath string) (string, error) {
	// var hasher = sha256.New()
	var hasher = crc32.NewIEEE()
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

// HashSHA256 computes the SHA256 hash of a fileops.
func HashSHA256(filePath string) (string, error) {
	var hasher = sha256.New()
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
