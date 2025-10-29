package cut

import (
	"testing"
)

func TestParseFields(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []int
		wantErr bool
	}{
		{
			name:    "single field",
			input:   "1",
			want:    []int{1},
			wantErr: false,
		},
		{
			name:    "multiple fields",
			input:   "1,3,5",
			want:    []int{1, 3, 5},
			wantErr: false,
		},
		{
			name:    "simple range",
			input:   "1-3",
			want:    []int{1, 2, 3},
			wantErr: false,
		},
		{
			name:    "range with spaces",
			input:   "1 - 3",
			want:    []int{1, 2, 3},
			wantErr: false,
		},
		{
			name:    "mixed fields and ranges",
			input:   "1,3-5,7",
			want:    []int{1, 3, 4, 5, 7},
			wantErr: false,
		},
		{
			name:    "multiple ranges",
			input:   "1-2,4-5",
			want:    []int{1, 2, 4, 5},
			wantErr: false,
		},
		{
			name:    "duplicate fields",
			input:   "1,2,1",
			want:    []int{1, 2, 1},
			wantErr: false,
		},
		{
			name:    "fields with many spaces",
			input:   " 1 , 3 - 5 , 7 ",
			want:    []int{1, 3, 4, 5, 7},
			wantErr: false,
		},
		{
			name:    "overlapping ranges",
			input:   "1-3,2-4",
			want:    []int{1, 2, 3, 2, 3, 4},
			wantErr: false,
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
		},
		{
			name:    "only commas",
			input:   ",,",
			wantErr: true,
		},
		{
			name:    "invalid range format - triple",
			input:   "1-2-3",
			wantErr: true,
		},
		{
			name:    "invalid lower bound",
			input:   "a-3",
			wantErr: true,
		},
		{
			name:    "invalid upper bound",
			input:   "1-b",
			wantErr: true,
		},
		{
			name:    "negative field",
			input:   "-1",
			wantErr: true,
		},
		{
			name:    "zero field",
			input:   "0",
			wantErr: true,
		},
		{
			name:    "negative in range",
			input:   "-1-3",
			wantErr: true,
		},
		{
			name:    "zero in range",
			input:   "0-3",
			wantErr: true,
		},
		{
			name:    "reverse range",
			input:   "5-3",
			wantErr: true,
		},
		{
			name:    "empty field in middle",
			input:   "1,,3",
			wantErr: true,
		},
		{
			name:    "only dash",
			input:   "-",
			wantErr: true,
		},
		{
			name:    "dash at start",
			input:   "-3",
			wantErr: true,
		},
		{
			name:    "dash at end",
			input:   "3-",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseFields(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseError: expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseError: %v", err)
				return
			}

			if !equalSlices(result, tt.want) {
				t.Errorf("ParseError: Input: %v; Got: %v; Want: %v", tt.input, result, tt.want)
			}
		})
	}

}

func equalSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
