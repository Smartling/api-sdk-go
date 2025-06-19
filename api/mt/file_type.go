package mt

//go:generate stringer -type=FileType
type FileType int

const (
	DOCX FileType = iota + 1
	DOCM
	RTF
	PPTX
	XLSX
	IDML
	RESX
	PLAIN_TEXT
	XML
	HTML
	PRES
	SRT
	MARKDOWN
	DITA
	VTT
	FLARE
	SVG
	XLIFF2
	CSV
	JSON
	XLSX_TEMPLATE
)
