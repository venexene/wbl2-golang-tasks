package unpack

import (
	"testing"
)

func TestUnpackString(t *testing.T) {
	tests := []struct {
		name string
		input string
		want string 
		wantErr bool
	} {
		{
            name:    "simple sequence",
            input:   "a4bc2d5e",
            want:    "aaaabccddddde",
            wantErr: false,
        },
        {
            name:    "single letters",
            input:   "abcd",
            want:    "abcd",
            wantErr: false,
        },
        {
            name:    "empty string",
            input:   "",
            want:    "",
            wantErr: false,
        },
		 {
            name:    "only digits",
            input:   `45`,
            want:    "",
            wantErr: true,
        },
        {
            name:    "escape digits",
            input:   `qwe\4\5`,
            want:    "qwe45",
            wantErr: false,
        },
        {
            name:    "escape letters",
            input:   `qwe\45`,
            want:    "qwe44444",
            wantErr: false,
        },
        {
            name:    "escape backslash",
            input:   `qwe\\5`,
            want:    `qwe\\\\\`,
            wantErr: false,
        },
        {
            name:    "starts with digit",
            input:   "4a5b6",
            want:    "",
            wantErr: true,
        },
        {
            name:    "invalid escape",
            input:   `qwe\`,
            want:    "",
            wantErr: true,
        },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unpack(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unpack error: Err: %v; WantErr: %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Unpack error: Got: %v; Want: %v", got, tt.want)
			}
		}) 
	}
}


func TestLexer(t *testing.T) {
	tests := []struct{
		input string
		want []Token
	} {
		 {
            input: "a2b3c4d5",
            want: []Token{
                {Symbol, "a"},
                {Coefficient, "2"},
                {Symbol, "b"},
                {Coefficient, "3"},
				{Symbol, "c"},
                {Coefficient, "4"},
				{Symbol, "d"},
                {Coefficient, "5"},
            },
        },
        {
            input: `a\3b\\e`,
            want: []Token{
                {Symbol, "a"},
                {Symbol, "3"},
                {Symbol, "b"},
				{Symbol, "\\"},
				{Symbol, "e"},
            },
        },
	}

	for _, tt := range tests {
		tokens, err := tokenize(tt.input)
        if err != nil {
            t.Errorf("Tokenize Error: Input:%v; Failed: %v", tt.input, err)
            continue
        }

        if len(tokens) != len(tt.want) {
            t.Errorf("Tokenize Error: Input:%v; Got Length: %d; Want Length: %d", tt.input, len(tokens), len(tt.want))
            continue
        }

        for i := range tokens {
            if tokens[i] != tt.want[i] {
                t.Errorf("Tokenize Error: Input: %v; Got: %v; Want %v", tt.input, tokens[i], tt.want[i])
            }
        }
	}
}