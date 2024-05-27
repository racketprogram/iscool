// main_test.go
package main

import (
	"bytes"
	"io"
	"os"
	"regexp"
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

	if len(expectedLines) != len(actualLines) {
		return false
	}

	timePattern := `\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`
	re := regexp.MustCompile(timePattern)

	for i := range expectedLines {
		expectedLine := re.ReplaceAllString(expectedLines[i], "")
		actualLine := re.ReplaceAllString(actualLines[i], "")

		if strings.TrimSpace(expectedLine) != strings.TrimSpace(actualLine) {
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
		// list folder and file
		{"register", []string{"user"}, "Add user successfully.\n"},
		{"create-folder", []string{"user", "folder_a"}, "Create folder_a successfully.\n"},
		{"list-folders", []string{"user"}, "folder_a 2000-01-01 20:34:19 user\n"},
		{"create-folder", []string{"user", "folder_b"}, "Create folder_b successfully.\n"},
		{"list-folders", []string{"user"}, "folder_a 2000-01-01 20:34:19 user\nfolder_b 2000-01-01 20:34:19 user\n"},
		{"list-folders", []string{"user"}, "folder_a 2000-01-01 20:34:19 user\nfolder_b 2000-01-01 20:34:19 user\n"},
		{"list-folders", []string{"user", "--sort-name", "asc"}, "folder_a 2000-01-01 20:34:19 user\nfolder_b 2000-01-01 20:34:19 user\n"},
		{"list-folders", []string{"user", "--sort-name", "desc"}, "folder_b 2000-01-01 20:34:19 user\nfolder_a 2000-01-01 20:34:19 user\n"},
		{"list-folders", []string{"user", "--sort-created", "asc"}, "folder_a 2000-01-01 20:34:19 user\nfolder_b 2000-01-01 20:34:19 user\n"},
		{"list-folders", []string{"user", "--sort-created", "desc"}, "folder_b 2000-01-01 20:34:19 user\nfolder_a 2000-01-01 20:34:19 user\n"},
		{"list-folders", []string{"user", "--sort-created"}, "Usage: list-folders [username] [--sort-name|--sort-created] [asc|desc]\n"},
		{"create-file", []string{"user", "folder_a", "file_a"}, "Create file_a in user/folder_a successfully.\n"},
		{"create-file", []string{"user", "folder_a", "file_b"}, "Create file_b in user/folder_a successfully.\n"},
		{"list-files", []string{"user", "folder_a"}, "file_a 2000-01-01 20:34:19 folder_a user\nfile_b 2000-01-01 20:34:19 folder_a user\n"},
		{"list-files", []string{"user", "folder_a", "--sort-name", "asc"}, "file_a 2000-01-01 20:34:19 folder_a user\nfile_b 2000-01-01 20:34:19 folder_a user\n"},
		{"list-files", []string{"user", "folder_a", "--sort-name", "desc"}, "file_b 2000-01-01 20:34:19 folder_a user\nfile_a 2000-01-01 20:34:19 folder_a user\n"},
		{"list-files", []string{"user", "folder_a", "--sort-created", "asc"}, "file_a 2000-01-01 20:34:19 folder_a user\nfile_b 2000-01-01 20:34:19 folder_a user\n"},
		{"list-files", []string{"user", "folder_a", "--sort-created", "desc"}, "file_b 2000-01-01 20:34:19 folder_a user\nfile_a 2000-01-01 20:34:19 folder_a user\n"},

		// quote name
		{"register", []string{"user 0"}, "Add \"user 0\" successfully.\n"},
		{"create-folder", []string{"user 0", "folder b", "folder b description"}, "Create \"folder b\" successfully.\n"},
		{"create-folder", []string{"user 0", "folder c", "folder c description"}, "Create \"folder c\" successfully.\n"},
		{"list-folders", []string{"user 0"}, "\"folder b\" \"folder b description\" 2000-01-01 20:34:19 \"user 0\"\n\"folder c\" \"folder c description\" 2000-01-01 20:34:19 \"user 0\"\n"},
		{"create-file", []string{"user 0", "folder c", "file c-1"}, "Create \"file c-1\" in \"user 0\"/\"folder c\" successfully.\n"},
		{"create-file", []string{"user 0", "folder c", "file c-2"}, "Create \"file c-2\" in \"user 0\"/\"folder c\" successfully.\n"},
		{"list-files", []string{"user 0", "folder c"}, "\"file c-1\" 2000-01-01 20:34:19 \"folder c\" \"user 0\"\n\"file c-2\" 2000-01-01 20:34:19 \"folder c\" \"user 0\"\n"},

		// delete file and folder
		{"register", []string{"user1"}, "Add user1 successfully.\n"},
		{"create-folder", []string{"user1", "folder1"}, "Create folder1 successfully.\n"},
		{"list-folders", []string{"user1"}, "folder1 2000-01-01 20:34:19 user1\n"},
		{"create-file", []string{"user1", "folder1", "file1"}, "Create file1 in user1/folder1 successfully.\n"},
		{"list-files", []string{"user1", "folder1"}, "file1 2000-01-01 20:34:19 folder1 user1\n"},
		{"delete-file", []string{"user1", "folder1", "file1"}, "Delete file1 in user1/folder1 successfully.\n"},
		{"delete-folder", []string{"user1", "folder1"}, "Delete folder1 successfully.\n"},
		{"list-folders", []string{"user1"}, "Warning: The user1 doesn't have any folders.\n"},

		// already existed
		{"register", []string{"user2"}, "Add user2 successfully.\n"},
		{"register", []string{"user2"}, "Error: The user2 has already existed.\n"},
		{"create-folder", []string{"user2", "folder2"}, "Create folder2 successfully.\n"},
		{"create-folder", []string{"user2", "folder2"}, "Error: The folder2 has already existed.\n"},
		{"create-file", []string{"user2", "folder2", "file2"}, "Create file2 in user2/folder2 successfully.\n"},
		{"create-file", []string{"user2", "folder2", "file2"}, "Error: The file2 has already existed.\n"},

		// invalid name
		{"register", []string{"!"}, "Error: The ! contains invalid chars.\n"},
		{"register", []string{"user3"}, "Add user3 successfully.\n"},
		{"create-folder", []string{"user3", "!"}, "Error: The ! contains invalid chars.\n"},
		{"create-folder", []string{"user3", "folder3"}, "Create folder3 successfully.\n"},
		{"create-file", []string{"user3", "folder3", "! !"}, "Error: The \"! !\" contains invalid chars.\n"},

		// user doesn't exist
		{"create-folder", []string{"user 4", "a"}, "Error: The \"user 4\" doesn't exist.\n"},

		// rename folder
		{"register", []string{"user5"}, "Add user5 successfully.\n"},
		{"create-folder", []string{"user5", "folder5"}, "Create folder5 successfully.\n"},
		{"create-file", []string{"user5", "folder5", "file5"}, "Create file5 in user5/folder5 successfully.\n"},
		{"rename-folder", []string{"user5", "folder5", "folder5-1"}, "Rename folder5 to folder5-1 successfully.\n"},
		{"list-folders", []string{"user5"}, "folder5-1 2000-01-01 20:34:19 user5\n"},
		{"list-files", []string{"user5", "folder5-1"}, "file5 2000-01-01 20:34:19 folder5-1 user5\n"},

		// case insensitive and description
		{"register", []string{"UseR6"}, "Add user6 successfully.\n"},
		{"create-folder", []string{"UseR6", "FoldeR6", "FoldeR6 DESCRIPTION!"}, "Create folder6 successfully.\n"},
		{"create-file", []string{"UseR6", "FoldeR6", "FilE6", "FilE6 DESCRIPTION!"}, "Create file6 in user6/folder6 successfully.\n"},
		{"list-folders", []string{"UseR6"}, "folder6 \"FoldeR6 DESCRIPTION!\" 2000-01-01 20:34:19 user6\n"},
		{"list-files", []string{"UseR6", "FoldeR6"}, "file6 \"FilE6 DESCRIPTION!\" 2000-01-01 20:34:19 folder6 user6\n"},
	}

	for _, tt := range tests {
		t.Run(tt.command, func(t *testing.T) {
			output := captureOutput(func() {
				handleCommand(tt.command, tt.args)
			})
			if !checkOutput(tt.expected, output) {
				t.Errorf("command: %v, args: %v\nexpected: %q\nbut got: %q", tt.command, tt.args, tt.expected, output)
			}
		})
	}
}
