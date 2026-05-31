package glossary

type ImportGlossaryRequest struct {
	File        []byte
	FileName    string
	MediaType   string
	ArchiveMode bool
}

type ImportGlossaryResponse struct {
	Code               int
	GlossaryUID        string
	ImportUID          string
	ImportStatus       string
	EntryChanges       ImportEntryChanges
	TranslationChanges []ImportTranslationChanges
	Warnings           []ImportWarning
}

type ImportEntryChanges struct {
	NewEntries           int
	ExistingEntryUpdates int
	NotMatchedEntries    int
	EntriesToArchive     int
}

type ImportTranslationChanges struct {
	LocaleID             string
	NewTranslations      int
	UpdatedTranslations  int
	TranslationsToRemove int
}

type ImportWarning struct {
	Key     string
	Message string
}

func toImportGlossaryResponse(r importGlossary, code int) ImportGlossaryResponse {
	res := ImportGlossaryResponse{
		Code:         code,
		GlossaryUID:  r.Response.Data.GlossaryImport.GlossaryUid,
		ImportUID:    r.Response.Data.GlossaryImport.ImportUid,
		ImportStatus: r.Response.Data.GlossaryImport.ImportStatus,
		EntryChanges: ImportEntryChanges{
			NewEntries:           r.Response.Data.EntryChanges.NewEntries,
			ExistingEntryUpdates: r.Response.Data.EntryChanges.ExistingEntryUpdates,
			NotMatchedEntries:    r.Response.Data.EntryChanges.NotMatchedEntries,
			EntriesToArchive:     r.Response.Data.EntryChanges.EntriesToArchive,
		},
	}
	res.TranslationChanges = make([]ImportTranslationChanges, len(r.Response.Data.TranslationChanges))
	for i, t := range r.Response.Data.TranslationChanges {
		res.TranslationChanges[i] = ImportTranslationChanges{
			LocaleID:             t.LocaleId,
			NewTranslations:      t.NewTranslations,
			UpdatedTranslations:  t.UpdatedTranslations,
			TranslationsToRemove: t.TranslationsToRemove,
		}
	}
	res.Warnings = make([]ImportWarning, len(r.Response.Data.Warnings))
	for i, w := range r.Response.Data.Warnings {
		res.Warnings[i] = ImportWarning{Key: w.Key, Message: w.Message}
	}
	return res
}

type importGlossary struct {
	Response struct {
		Code string `json:"code"`
		Data struct {
			GlossaryImport struct {
				GlossaryUid  string `json:"glossaryUid"`
				ImportUid    string `json:"importUid"`
				ImportStatus string `json:"importStatus"`
			} `json:"glossaryImport"`
			EntryChanges struct {
				NewEntries           int `json:"newEntries"`
				ExistingEntryUpdates int `json:"existingEntryUpdates"`
				NotMatchedEntries    int `json:"notMatchedEntries"`
				EntriesToArchive     int `json:"entriesToArchive"`
			} `json:"entryChanges"`
			TranslationChanges []struct {
				LocaleId             string `json:"localeId"`
				NewTranslations      int    `json:"newTranslations"`
				UpdatedTranslations  int    `json:"updatedTranslations"`
				TranslationsToRemove int    `json:"translationsToRemove"`
			} `json:"translationChanges"`
			Warnings []struct {
				Key     string `json:"key"`
				Message string `json:"message"`
			} `json:"warnings"`
		} `json:"data"`
	} `json:"response"`
}
