package glossary

import (
	"io"
	"time"
)

type ExportGlossaryRequest struct {
	Format        string               `json:"format"`
	TbxVersion    string               `json:"tbxVersion,omitempty"`
	Filter        ExportGlossaryFilter `json:"filter"`
	FocusLocaleId string               `json:"focusLocaleId,omitempty"`
	LocaleIds     []string             `json:"localeIds"`
	SkipEntries   bool                 `json:"skipEntries,omitempty"`
}

type ExportGlossaryFilter struct {
	Query                      string                      `json:"query,omitempty"`
	LocaleIds                  []string                    `json:"localeIds,omitempty"`
	EntryUids                  []string                    `json:"entryUids,omitempty"`
	EntryState                 string                      `json:"entryState,omitempty"`
	MissingTranslationLocaleId string                      `json:"missingTranslationLocaleId,omitempty"`
	PresentTranslationLocaleId string                      `json:"presentTranslationLocaleId,omitempty"`
	DntLocaleId                string                      `json:"dntLocaleId,omitempty"`
	ReturnFallbackTranslations bool                        `json:"returnFallbackTranslations,omitempty"`
	Labels                     *ExportGlossaryLabelsFilter `json:"labels,omitempty"`
	DntTermSet                 bool                        `json:"dntTermSet,omitempty"`
	Created                    *ExportGlossaryDateFilter   `json:"created,omitempty"`
	LastModified               *ExportGlossaryDateFilter   `json:"lastModified,omitempty"`
	CreatedBy                  *ExportGlossaryUserFilter   `json:"createdBy,omitempty"`
	LastModifiedBy             *ExportGlossaryUserFilter   `json:"lastModifiedBy,omitempty"`
	Paging                     ExportGlossaryPaging        `json:"paging"`
	Sorting                    *ExportGlossarySorting      `json:"sorting,omitempty"`
}

type ExportGlossaryLabelsFilter struct {
	Type string `json:"type,omitempty"`
}

type ExportGlossaryDateFilter struct {
	Level string    `json:"level,omitempty"`
	Date  time.Time `json:"date"`
	Type  string    `json:"type,omitempty"`
}

type ExportGlossaryUserFilter struct {
	Level   string   `json:"level,omitempty"`
	UserIds []string `json:"userIds,omitempty"`
}

type ExportGlossaryPaging struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type ExportGlossarySorting struct {
	Field     string `json:"field,omitempty"`
	Direction string `json:"direction,omitempty"`
	LocaleId  string `json:"localeId,omitempty"`
}

type ExportGlossaryResponse struct {
	Code          int
	Filename      string
	ContentType   string
	ContentLength int64
	Data          io.ReadCloser
}
