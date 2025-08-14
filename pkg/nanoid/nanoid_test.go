package nanoid

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		alphabet string
		size     int
		wantErr  bool
	}{
		{"abcdef", 10, false},
		{"", 10, true},
		{"abcdef", 0, true},
		{"abcdef", -1, true},
		{"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 21, false},
	}

	for _, tt := range tests {
		t.Run(tt.alphabet, func(t *testing.T) {
			_, err := Generate(tt.alphabet, tt.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMustGenerate(t *testing.T) {
	tests := []struct {
		alphabet  string
		size      int
		wantPanic bool
	}{
		{"abcdef", 10, false},
		{"", 10, true},
		{"abcdef", 0, true},
		{"abcdef", -1, true},
		{"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 21, false},
	}

	for _, tt := range tests {
		t.Run(tt.alphabet, func(t *testing.T) {
			defer func() {
				if r := recover(); (r != nil) != tt.wantPanic {
					t.Errorf("MustGenerate() panic = %v, wantPanic %v", r, tt.wantPanic)
				}
			}()
			MustGenerate(tt.alphabet, tt.size)
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		lengths []int
		wantErr bool
	}{
		{nil, false},
		{[]int{10}, false},
		{[]int{-1}, true},
		{[]int{10, 20}, true},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			_, err := New(tt.lengths...)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMust(t *testing.T) {
	tests := []struct {
		lengths   []int
		wantPanic bool
	}{
		{nil, false},
		{[]int{10}, false},
		{[]int{-1}, true},
		{[]int{10, 20}, true},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			defer func() {
				if r := recover(); (r != nil) != tt.wantPanic {
					t.Errorf("Must() panic = %v, wantPanic %v", r, tt.wantPanic)
				}
			}()
			Must(tt.lengths...)
		})
	}
}
