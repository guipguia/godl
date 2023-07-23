package tests

// import "testing"

// func TestMultiOSPaths(t *testing.T) {
// 	tests := []struct {
// 		input    string
// 		expected string
// 	}{
// 		{"/Downloads", "/home/test_user/Downloads"},                         // Linux-style path
// 		{"C:\\Users\\Test User\\Downloads", "D:\\Downloads"},                // Windows-style path
// 		{"C:/Users/Test User/Downloads", "C:\\Users\\Test User\\Downloads"}, // Mixed-style path (works on Windows)
// 	}

// 	for _, test := range tests {
// 		result := NormalizePath(test.input)
// 		if result != test.expected {
// 			t.Errorf("NormalizePath(%q) = %q, expected %q", test.input, result, test.expected)
// 		}
// 	}
// }
