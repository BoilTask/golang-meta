package metafile

import (
	"testing"
)

func TestGetFileNameWithoutExt(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{
			path:     "meta/meta-file/library.go",
			expected: "library",
		}}
	for _, test := range tests {
		if output := GetFileNameWithoutExt(test.path); output != test.expected {
			t.Errorf("Test failed: %v inputted, %v expected, %v received", test.path, test.expected, output)
		}
	}
}
