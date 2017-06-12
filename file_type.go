package smartling

// FileType represents file type format used in Smartling API.
type FileType string

// Android and next are types that are supported by Smartling API.
const (
	Android        FileType = "android"
	IOS            FileType = "ios"
	Gettext        FileType = "gettext"
	HTML           FileType = "html"
	JavaProperties FileType = "javaProperties"
	YAML           FileType = "yaml"
	XLIFF          FileType = "xliff"
	XML            FileType = "xml"
	JSON           FileType = "json"
	DOCX           FileType = "docx"
	PPTX           FileType = "pptx"
	XLSX           FileType = "xlsx"
	IDML           FileType = "idml"
	Qt             FileType = "qt"
	Resx           FileType = "resx"
	Plaintext      FileType = "plaintext"
	CSV            FileType = "csv"
	Stringsdict    FileType = "stringsdict"
)
