// main_test.go
package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"virtual-file-system/internal"
)

func captureOutput(f func()) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	stdout := os.Stdout
	stderr := os.Stderr
	os.Stdout = writer
	os.Stderr = writer

	out := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, reader)
		out <- buf.String()
	}()

	f()

	writer.Close()
	os.Stdout = stdout
	os.Stderr = stderr

	return <-out
}

func checkOutput(expected, actual string) bool {
	expectedLines := strings.Split(expected, "\n")
	actualLines := strings.Split(actual, "\n")

	if len(expectedLines) > len(actualLines) {
		return false
	}

	for i := range expectedLines {
		if !strings.HasPrefix(actualLines[i], expectedLines[i]) {
			return false
		}
	}
	return true
}

func TestHandleCommand(t *testing.T) {
	// Set up mock data
	mockUsers := make(map[string]*internal.User)
	internal.UseMockData(mockUsers)

	tests := []struct {
		command  string
		args     []string
		expected string
	}{
		{"register", []string{"user"}, "Add user successfully."},
		{"create-folder", []string{"user", "folderA"}, "Create folderA successfully."},
		{"list-folders", []string{"user"}, "folderA"},
		{"create-folder", []string{"user", "folderB"}, "Create folderB successfully."},
		{"list-folders", []string{"user"}, "folderA\nfolderB"},
		{"list-folders", []string{"user", "--sort-name", "asc"}, "folderA\nfolderB"},
		{"list-folders", []string{"user", "--sort-name", "desc"}, "folderB\nfolderA"},
		{"list-folders", []string{"user", "--sort-created", "asc"}, "folderA\nfolderB"},
		{"list-folders", []string{"user", "--sort-created", "desc"}, "folderB\nfolderA"},

		{"register", []string{"user a"}, "Add \"user a\" successfully."},
		{"create-folder", []string{"user a", "folder b", "folder b description"}, "Create \"folder b\" successfully."},
		{"list-folders", []string{"user a"}, "\"folder b\" \"folder b description\""},
		{"create-folder", []string{"user a", "folder c", "folder c description"}, "Create \"folder c\" successfully."},
		{"list-folders", []string{"user a"}, "\"folder b\" \"folder b description\"\n\"folder c\" \"folder c description\""},
		{"create-file", []string{"user a", "folder c", "file c"}, "Create \"file c\" in \"user a\"/\"folder c\" successfully."},

		{"register", []string{"user1"}, "Add user1 successfully.\n"},
		{"create-folder", []string{"user1", "folder1"}, "Create folder1 successfully."},
		{"list-folders", []string{"user1"}, "folder1"},
		{"delete-folder", []string{"user1", "folder1"}, "Delete folder1 successfully."},
		{"list-folders", []string{"user1"}, ""},

		{"register", []string{"user2"}, "Add user2 successfully."},
		{"register", []string{"user2"}, "Error: The user2 has already existed."},
		{"create-folder", []string{"user2", "folder2"}, "Create folder2 successfully."},
		{"create-folder", []string{"user2", "folder2"}, "Error: The folder2 has already existed."},
		{"create-file", []string{"user2", "folder2", "file2"}, "Create file2 in user2/folder2 successfully."},
		{"create-file", []string{"user2", "folder2", "file2"}, "Error: The file2 has already existed."},

		{"register", []string{"!"}, "Error: The ! contains invalid chars."},
		{"register", []string{"user3"}, "Add user3 successfully.\n"},
		{"create-folder", []string{"user3", "!"}, "Error: The ! contains invalid chars."},
		{"create-folder", []string{"user3", "folder3"}, "Create folder3 successfully."},
		{"create-file", []string{"user3", "folder3", "! !"}, "Error: The \"! !\" contains invalid chars."},

		{"create-folder", []string{"user 4", "a"}, "Error: The \"user 4\" doesn't exist."},
	}

	for _, tt := range tests {
		t.Run(tt.command, func(t *testing.T) {
			output := captureOutput(func() {
				handleCommand(tt.command, tt.args)
			})
			if !checkOutput(tt.expected, output) {
				t.Errorf("expected %q but got %q", tt.expected, output)
			}
		})
	}
}
