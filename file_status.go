package smartling

// FileStatus represents current file status in the Smartling system.
type FileStatus struct {
	// FileURI is a unique path to file in Smartling system.
	FileURI string

	// LastUploaded refers to time when file was uploaded.
	LastUploaded UTC

	// FileType is a file type identifier.
	FileType FileType

	// HasInstructions specifies does files have instructions or not.
	HasInstructions bool
}
