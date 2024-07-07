package storage

import (
	"encoding/json"
	"io/ioutil"

	"github.com/draychev/file-catalog/pkg/fileops"
)

// SerializeFileMeta serializes the fileops metadata to a JSON fileops.
func SerializeFileMeta(metas []fileops.FileMeta, outputPath string) error {
	data, err := json.MarshalIndent(metas, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(outputPath, data, 0644)
}
