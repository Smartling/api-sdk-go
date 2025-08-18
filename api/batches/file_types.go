package batches

// Type is type for Batch file types
//
//go:generate stringer -type=Type -output=file_types_string.go
type Type int

const (
	JAVA_PROPERTIES Type = iota + 1
	IOS
	STRINGSDICT
	ANDROID
	GETTEXT
	PHP_RESOURCE
	RESX
	XLIFF
	YAML
	JSON
	XML
	HTML
	FREEMARKER
	DOCX
	PPTX
	XLSX
	IDML
	XLS
	DOC
	QT
	CSV
	TMX
	PLAIN_TEXT
	PPT
	PRES
	MADCAP
	SRT
	MARKDOWN
	DITA
	DITA_ZIP
	VTT
	PDF
	RTF
	FLARE
	FLUENT
	SVG
	DOCM
	ARB
	INDD
	XCSTRINGS
	VSDX
	VSDM
)

var TypeByExt = map[string]Type{
	".properties":  JAVA_PROPERTIES,
	".strings":     IOS,
	".stringsdict": STRINGSDICT,
	".po":          GETTEXT,
	".pot":         GETTEXT,
	".php":         PHP_RESOURCE,
	".resx":        RESX,
	".resw":        RESX,
	".xlf":         XLIFF,
	".xliff":       XLIFF,
	".yml":         YAML,
	".yaml":        YAML,
	".json":        JSON,
	".js":          JSON,
	".xml":         XML,
	".html":        HTML,
	".htm":         HTML,
	".docx":        DOCX,
	".pptx":        PPTX,
	".xlsx":        XLSX,
	".idml":        IDML,
	".xls":         XLS,
	".doc":         DOC,
	".ts":          QT,
	".csv":         CSV,
	".tmx":         TMX,
	".txt":         PLAIN_TEXT,
	".ppt":         PPT,
	".pres":        PRES,
	".srt":         SRT,
	".markdown":    MARKDOWN,
	".md":          MARKDOWN,
	".dita":        DITA,
	".ditamap":     DITA,
	".zip":         DITA_ZIP,
	".vtt":         VTT,
	".pdf":         PDF,
	".rtf":         RTF,
	".flprjzip":    FLARE,
	".ftl":         FLUENT,
	".svg":         SVG,
	".docm":        DOCM,
	".arb":         ARB,
	".indd":        INDD,
	".xcstrings":   XCSTRINGS,
	".vsdx":        VSDX,
	".vsdm":        VSDM,
}
