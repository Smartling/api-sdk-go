package file

// Type is type for file types
//
//go:generate stringer -type=Type
type Type int

const (
	DOCX Type = iota + 1
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

// TypeByExt contains map with FileType by file extension
var TypeByExt = map[string]Type{
	".docx":          DOCX,
	".docm":          DOCM,
	".rtf":           RTF,
	".pptx":          PPTX,
	".xlsx":          XLSX,
	".idml":          IDML,
	".resx":          RESX,
	".txt":           PLAIN_TEXT,
	".xml":           XML,
	".html":          HTML,
	".htm":           HTML,
	".pres":          PRES,
	".srt":           SRT,
	".md":            MARKDOWN,
	".markdown":      MARKDOWN,
	".dita":          DITA,
	".vtt":           VTT,
	".zip":           FLARE,
	".svg":           SVG,
	".xlf":           XLIFF2,
	".xliff":         XLIFF2,
	".csv":           CSV,
	".json":          JSON,
	".xlsx_template": XLSX_TEMPLATE,
}
