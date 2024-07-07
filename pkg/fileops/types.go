package fileops

// FileMeta holds metadata of a fileops.
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
