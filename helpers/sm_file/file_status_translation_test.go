package smfile

import "testing"

func TestFileStatusTranslation_ProgressPercent(t *testing.T) {
	tests := []struct {
		name                string
		given               FileStatusTranslation
		totalStringCount    int
		wantProgressPercent int
		wantErr             bool
	}{
		{
			name:             "progress = 0%, 4 strings are ready for translation but translator didn't start yet",
			totalStringCount: 10,
			given: FileStatusTranslation{
				AuthorizedStringCount: 4,
				CompletedStringCount:  0,
				ExcludedStringCount:   1,
			},
			wantProgressPercent: 0,
		},
		{
			name:             "progress = 55%, it's common case",
			totalStringCount: 10,
			given: FileStatusTranslation{
				AuthorizedStringCount: 4,
				CompletedStringCount:  5,
				ExcludedStringCount:   1,
			},
			wantProgressPercent: 55,
		},
		{
			name:             "progress = 55%, 2 strings are waiting for translation and 2 more for authorization",
			totalStringCount: 10,
			given: FileStatusTranslation{
				AuthorizedStringCount: 2,
				CompletedStringCount:  5,
				ExcludedStringCount:   1,
			},
			wantProgressPercent: 55,
		},
		{
			name:             "progress = 55%, still 4 strings are waiting for decision (authorize/exclude)",
			totalStringCount: 10,
			given: FileStatusTranslation{
				AuthorizedStringCount: 0,
				CompletedStringCount:  5,
				ExcludedStringCount:   1,
			},
			wantProgressPercent: 55,
		},
		{
			name:             "progress = 90%, file was uploaded without authorization, user will do this later or will add file to job",
			totalStringCount: 10,
			given: FileStatusTranslation{
				AuthorizedStringCount: 0,
				CompletedStringCount:  9,
				ExcludedStringCount:   0,
			},
			wantProgressPercent: 90,
		},
		{
			name:             "progress = 0%, file was uploaded without authorization, user will do this later or will add file to job.",
			totalStringCount: 10,
			given: FileStatusTranslation{
				AuthorizedStringCount: 0,
				CompletedStringCount:  0,
				ExcludedStringCount:   0,
			},
			wantProgressPercent: 0,
		},
		{
			name:             "progress = 100%, nothing to translate, user excluded all content",
			totalStringCount: 10,
			given: FileStatusTranslation{
				AuthorizedStringCount: 0,
				CompletedStringCount:  0,
				ExcludedStringCount:   10,
			},
			wantProgressPercent: 100,
		},
		{
			name:             "progress = 99%, must return 99% even if 99.9999% translated",
			totalStringCount: 1000000,
			given: FileStatusTranslation{
				AuthorizedStringCount: 0,
				CompletedStringCount:  999999,
				ExcludedStringCount:   0,
			},
			wantProgressPercent: 99,
		},
		{
			name:             "Log error, totalStringCount must be greater than 0",
			totalStringCount: 0,
			given: FileStatusTranslation{
				AuthorizedStringCount: 0,
				CompletedStringCount:  10,
				ExcludedStringCount:   0,
			},
			wantProgressPercent: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotProgressPercent, gotErr := tt.given.ProgressPercent(tt.totalStringCount)
			if gotProgressPercent != tt.wantProgressPercent {
				t.Errorf("ProgressPercent() = %v, wantProgressPercent %v", gotProgressPercent, tt.wantProgressPercent)
			}
			if tt.wantErr != (gotErr != nil) {
				t.Errorf("err = %v, wantProgressPercent %v", gotProgressPercent, tt.wantProgressPercent)
			}
		})
	}
}
