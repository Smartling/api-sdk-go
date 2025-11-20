package smfile

import "testing"

func TestFileStatusTranslation_ProgressPercent(t *testing.T) {
	tests := []struct {
		name  string
		given FileStatusTranslation
		want  int
	}{
		{
			name:  "Zero total",
			given: FileStatusTranslation{CompletedStringCount: 0, AuthorizedStringCount: 0},
		},
		{
			name:  "50 percent",
			given: FileStatusTranslation{CompletedStringCount: 100, AuthorizedStringCount: 100},
			want:  50,
		},
		{
			name:  "100 percent complete without authorized",
			given: FileStatusTranslation{CompletedStringCount: 50, AuthorizedStringCount: 0},
			want:  100,
		},
		{
			name:  "0 percent complete with only authorized",
			given: FileStatusTranslation{CompletedStringCount: 0, AuthorizedStringCount: 100},
			want:  0,
		},
		{
			name:  "33 percent (rounding down)",
			given: FileStatusTranslation{CompletedStringCount: 10, AuthorizedStringCount: 20},
			want:  33,
		},
		{
			name:  "67 percent (rounding up)",
			given: FileStatusTranslation{CompletedStringCount: 20, AuthorizedStringCount: 10},
			want:  66,
		},
		{
			name:  "Large numbers 75 percent",
			given: FileStatusTranslation{CompletedStringCount: 75000, AuthorizedStringCount: 25000},
			want:  75,
		},
		{
			name:  "Near zero",
			given: FileStatusTranslation{CompletedStringCount: 1, AuthorizedStringCount: 9},
			want:  10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.given.ProgressPercent(); got != tt.want {
				t.Errorf("ProgressPercent() = %v, want %v", got, tt.want)
			}
		})
	}
}
