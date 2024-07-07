package storage

import (
	"encoding/json"
	"io/ioutil"

	"github.com/draychev/file-catalog/pkg/fileops"
)

// DeserializeFileMeta deserializes the fileops metadata from a JSON fileops.
func DeserializeFileMeta(inputPath string) ([]fileops.FileMeta, error) {
	data, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return nil, err
	}

	var metas []fileops.FileMeta
	if err := json.Unmarshal(data, &metas); err != nil {
		return nil, err
	}

	return metas, nil
}
