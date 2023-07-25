package tests

import (
	"testing"

	"github.com/guipguia/godl/internal/util"
)

func TestBaseName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"/Downloads/file.txt", "file.txt"},                       // Linux-style path
		{"https://example.com/file.txt", "file.txt"},              //URL like
		{"C:\\Users\\Test User\\Downloads\\file.txt", "file.txt"}, // Windows-style path
		{"C:/Users/Test User/Downloads/file.txt", "file.txt"},     // Mixed-style path (works on Windows)
	}

	for _, test := range tests {
		result := util.GetBaseName(test.input)
		if result != test.expected {
			t.Errorf("BaseName(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}
