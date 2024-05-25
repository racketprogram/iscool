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
		command   string
		args      []string
		expected  string
		shouldErr bool
	}{
		{"register", []string{"user"}, "Add user successfully.\n", false},
		{"create-folder", []string{"user", "folderA"}, "Create folderA successfully.\n", false},
		{"list-folders", []string{"user"}, "folderA", false},
		{"create-folder", []string{"user", "folderB"}, "Create folderB successfully.\n", false},
		{"list-folders", []string{"user"}, "folderA\nfolderB", false},
		{"list-folders", []string{"user", "--sort-name", "asc"}, "folderA\nfolderB", false},
		{"list-folders", []string{"user", "--sort-name", "desc"}, "folderB\nfolderA", false},
		{"list-folders", []string{"user", "--sort-created", "asc"}, "folderA\nfolderB", false},
		{"list-folders", []string{"user", "--sort-created", "desc"}, "folderB\nfolderA", false},

		{"register", []string{"user a"}, "Add \"user a\" successfully.\n", false},
		{"create-folder", []string{"user a", "folder b", "folder b description"}, "Create \"folder b\" successfully.\n", false},
		{"list-folders", []string{"user a"}, "\"folder b\" \"folder b description\"", false},
		{"create-folder", []string{"user a", "folder c", "folder c description"}, "Create \"folder c\" successfully.\n", false},
		{"list-folders", []string{"user a"}, "\"folder b\" \"folder b description\"\n\"folder c\" \"folder c description\"", false},

		{"register", []string{"u1"}, "Add u1 successfully.\n", false},
		{"create-folder", []string{"u1", "folderA"}, "Create folderA successfully.\n", false},
		{"list-folders", []string{"u1"}, "folderA", false},
		{"delete-folder", []string{"u1", "folderA"}, "Delete folderA successfully.\n", false},
		{"list-folders", []string{"u1"}, "", false},

		{"register", []string{"u2"}, "Add u2 successfully.\n", false},
		{"register", []string{"u2"}, "Error: The u2 has already existed.\n", false},

		// {"list-files", []string{"testuser", "testfolder"}, "", false},
		// {"delete-folder", []string{"testuser", "testfolder"}, "Delete folder \"testfolder\" successfully.\n", false},
		// {"delete-file", []string{"testuser", "testfolder", "testfile"}, "Delete file \"testfile\" successfully.\n", false},
		// {"rename-folder", []string{"testuser", "testfolder", "newfoldername"}, "Rename folder \"testfolder\" to \"newfoldername\" successfully.\n", false},
		// {"unknown", []string{}, "Unrecognized command\n", true},
	}

	for _, tt := range tests {
		t.Run(tt.command, func(t *testing.T) {
			output := captureOutput(func() {
				handleCommand(tt.command, tt.args)
			})
			if tt.shouldErr {
				if output == tt.expected {
					t.Errorf("expected an error but got none")
				}
			} else {
				if !checkOutput(tt.expected, output) {
					t.Errorf("expected %q but got %q", tt.expected, output)
				}
			}
		})
	}
}
