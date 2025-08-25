package smfile_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Smartling/api-sdk-go/api/batches"
	"github.com/Smartling/api-sdk-go/api/mt"
	smfile "github.com/Smartling/api-sdk-go/helpers/sm_file"
)

func TestParseType(t *testing.T) {
	type testCase[T fmt.Stringer] struct {
		name      string
		typeByExt map[string]T
		typ       string
		want      T
	}
	testsBatches := []testCase[batches.Type]{
		{
			name:      "YAML",
			typeByExt: batches.TypeByExt,
			typ:       "YAML",
			want:      batches.YAML,
		},
		{
			name:      "JAVA_PROPERTIES",
			typeByExt: batches.TypeByExt,
			typ:       "JAVA_PROPERTIES",
			want:      batches.JAVA_PROPERTIES,
		},
		{
			name:      "none",
			typeByExt: batches.TypeByExt,
			typ:       "TESTS",
			want:      0,
		},
	}
	for _, tt := range testsBatches {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := smfile.ParseType(batches.FirstType, batches.LastType, tt.typ); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseType() = %v, want %v", got, tt.want)
			}
		})
	}
	testsMT := []testCase[mt.Type]{
		{
			name:      "DOCX",
			typeByExt: mt.TypeByExt,
			typ:       "DOCX",
			want:      mt.DOCX,
		},
		{
			name:      "PLAIN_TEXT",
			typeByExt: mt.TypeByExt,
			typ:       "PLAIN_TEXT",
			want:      mt.PLAIN_TEXT,
		},
		{
			name:      "none",
			typeByExt: mt.TypeByExt,
			typ:       "TESTS",
			want:      0,
		},
	}
	for _, tt := range testsMT {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := smfile.ParseType(mt.FirstType, mt.LastType, tt.typ); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseType() = %v, want %v", got, tt.want)
			}
		})
	}
}
