package anagrams

import (
	"testing"
	"reflect"
)

func TestCreateFrequncyKey(t *testing.T) {
	tests := []struct {
		name string
		input string
		want string 
		wantErr bool
	} {
		{
            name:    "lower case 1",
            input:   "hello",
            want:    "e1h1l2o1",
            wantErr: false,
        },
		{
            name:    "lower case 2",
            input:   "test",
            want:    "e1s1t2",
            wantErr: false,
        },
        {
            name:    "upper case 1",
            input:   "Hello",
            want:    "e1h1l2o1",
            wantErr: false,
        },
        {
            name:    "upper case 2",
            input:   "HELLO",
            want:    "e1h1l2o1",
            wantErr: false,
        },
        {
            name:    "empty",
            input:   "",
            want:    "",
            wantErr: false,
        },
        {
            name:    "cyrillic",
            input:   "привет",
            want:    "в1е1и1п1р1т1",
            wantErr: true,
        },
        {
            name:    "uniqode",
            input:   "café",
            want:    "a1c1f1é1",
            wantErr: true,
        },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got:= createFrequencyKey(tt.input)
			if got != tt.want {
				t.Errorf("Frequency key error: Got: %v; Want: %v", got, tt.want)
			}
		}) 
	}
}

func TestFindAnagrams(t *testing.T) {
	tests := []struct {
		name string
		input []string
		want map[string][]string 
		wantErr bool
	} {
		{
            name:    "example",
            input:   []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик"},
            want:    map[string][]string{
						"пятак":  {"пятак", "пятка", "тяпка"},
						"листок": {"листок", "слиток", "столик"},
					},
            wantErr: false,
        },
		{
            name:    "mixed case",
            input:   []string{"Hello", "hello", "olelh", "World", "wrold"},
            want:    map[string][]string{
						"hello": {"hello", "olelh"},
						"world": {"world", "wrold"},
					},
            wantErr: false,
        },
        {
            name:    "no anagrams",
            input:   []string{"cat", "dog", "bird", "fish"},
            want:    map[string][]string{},
            wantErr: false,
        },
        {
            name:    "empty",
            input:   []string{},
            want:    map[string][]string{},
            wantErr: false,
        },
        {
            name:    "some test",
            input:   []string{"abcd", "bdac", "cabd", "def", "fed", "xyz", "yxz", "zxy"},
            want:    map[string][]string{
						"abcd": {"abcd", "bdac", "cabd"},
						"def": {"def", "fed"},
						"xyz": {"xyz", "yxz", "zxy"},
					},
            wantErr: false,
        },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got:= FindAnagrams(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		}) 
	}
}