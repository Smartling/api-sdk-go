package mt

import "fmt"

// Type is type for MT file types
//
//go:generate stringer -type=Type -output=file_types_string.go
type Type uint8

// MarshalText returns the stringer name so encoding/json (and any other
// encoding.TextMarshaler-aware encoder) emits the API-expected string form
// like "PLAIN_TEXT" instead of the raw uint8 value.
func (t Type) MarshalText() ([]byte, error) {
	if t < FirstType || t > LastType {
		return nil, fmt.Errorf("mt: invalid file Type value %d", t)
	}
	return []byte(t.String()), nil
}

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

	FirstType = DOCX
	LastType  = XLSX_TEMPLATE
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
