// internal/user_test.go
package internal

import (
	"testing"
)

func setupMockData() {
	mockUsers := make(map[string]*User)
	UseMockData(mockUsers)
}

func TestQuoteIfNeeded(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"hello world", "\"hello world\""},
		{"", ""},
		{" ", "\" \""},
	}

	for _, test := range tests {
		result := QuoteIfNeeded(test.input)
		if result != test.expected {
			t.Errorf("QuoteIfNeeded(%s) = %s; expected %s", test.input, result, test.expected)
		}
	}
}

func TestIsValidName(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"valid_name", true},
		{"valid name", true},
		{"valid-name", true},
		{"", false},
		{"a", true},
		{"a very long name that exceeds fifty characters in length should fail", false},
		{"invalid/name", false},
	}

	for _, test := range tests {
		result := isValidName(test.input)
		if result != test.expected {
			t.Errorf("isValidName(%s) = %v; expected %v", test.input, result, test.expected)
		}
	}
}

func TestRegisterUser(t *testing.T) {
	setupMockData()

	tests := []struct {
		username string
		expected error
	}{
		{"user1", nil},
		{"user1", errorAlreayExisted("user1")},
		{"invalid/user", errorInvalidChars("invalid/user")},
	}

	for _, test := range tests {
		err := RegisterUser(test.username)
		if err != nil && err.Error() != test.expected.Error() {
			t.Errorf("RegisterUser(%s) = %v; expected %v", test.username, err, test.expected)
		}
	}
}

func TestCreateFolder(t *testing.T) {
	setupMockData()
	RegisterUser("user1")

	tests := []struct {
		username    string
		foldername  string
		description string
		expected    error
	}{
		{"user1", "folder1", "desc1", nil},
		{"user1", "folder1", "desc1", errorAlreayExisted("folder1")},
		{"user1", "invalid/folder", "desc2", errorInvalidChars("invalid/folder")},
		{"user2", "folder2", "desc2", errorDoesntExisted("user2")},
	}

	for _, test := range tests {
		err := CreateFolder(test.username, test.foldername, test.description)
		if err != nil && err.Error() != test.expected.Error() {
			t.Errorf("CreateFolder(%s, %s, %s) = %v; expected %v", test.username, test.foldername, test.description, err, test.expected)
		}
	}
}

func TestCreateFile(t *testing.T) {
	setupMockData()
	RegisterUser("user1")
	CreateFolder("user1", "folder1", "desc1")

	tests := []struct {
		username    string
		foldername  string
		filename    string
		description string
		expected    error
	}{
		{"user1", "folder1", "file1", "desc1", nil},
		{"user1", "folder1", "file1", "desc1", errorAlreayExisted("file1")},
		{"user1", "folder1", "invalid/file", "desc2", errorInvalidChars("invalid/file")},
		{"user1", "folder2", "file2", "desc2", errorDoesntExisted("folder2")},
		{"user2", "folder1", "file2", "desc2", errorDoesntExisted("user2")},
	}

	for _, test := range tests {
		err := CreateFile(test.username, test.foldername, test.filename, test.description)
		if err != nil && err.Error() != test.expected.Error() {
			t.Errorf("CreateFile(%s, %s, %s, %s) = %v; expected %v", test.username, test.foldername, test.filename, test.description, err, test.expected)
		}
	}
}

func TestListFolders(t *testing.T) {
	setupMockData()
	RegisterUser("user1")
	CreateFolder("user1", "folderA", "descA")
	CreateFolder("user1", "folderB", "descB")

	tests := []struct {
		username string
		sortBy   string
		order    string
		expected []string
	}{
		{"user1", "name", "asc", []string{"folderA", "folderB"}},
		{"user1", "name", "desc", []string{"folderB", "folderA"}},
		{"user1", "created", "asc", []string{"folderA", "folderB"}},
		{"user1", "created", "desc", []string{"folderB", "folderA"}},
	}

	for _, test := range tests {
		folders, err := ListFolders(test.username, test.sortBy, test.order)
		if err != nil {
			t.Errorf("ListFolders(%s, %s, %s) returned error: %v", test.username, test.sortBy, test.order, err)
		}

		var folderNames []string
		for _, folder := range folders {
			folderNames = append(folderNames, folder.Name)
		}

		for i, name := range test.expected {
			if folderNames[i] != name {
				t.Errorf("ListFolders(%s, %s, %s) = %v; expected %v", test.username, test.sortBy, test.order, folderNames, test.expected)
				break
			}
		}
	}
}

func TestListFiles(t *testing.T) {
	setupMockData()
	RegisterUser("user1")
	CreateFolder("user1", "folder1", "desc1")
	CreateFile("user1", "folder1", "fileA", "descA")
	CreateFile("user1", "folder1", "fileB", "descB")

	tests := []struct {
		username   string
		foldername string
		sortBy     string
		order      string
		expected   []string
	}{
		{"user1", "folder1", "name", "asc", []string{"fileA", "fileB"}},
		{"user1", "folder1", "name", "desc", []string{"fileB", "fileA"}},
		{"user1", "folder1", "created", "asc", []string{"fileA", "fileB"}},
		{"user1", "folder1", "created", "desc", []string{"fileB", "fileA"}},
	}

	for _, test := range tests {
		files, err := ListFiles(test.username, test.foldername, test.sortBy, test.order)
		if err != nil {
			t.Errorf("ListFiles(%s, %s, %s, %s) returned error: %v", test.username, test.foldername, test.sortBy, test.order, err)
		}

		var fileNames []string
		for _, file := range files {
			fileNames = append(fileNames, file.Name)
		}

		for i, name := range test.expected {
			if fileNames[i] != name {
				t.Errorf("ListFiles(%s, %s, %s, %s) = %v; expected %v", test.username, test.foldername, test.sortBy, test.order, fileNames, test.expected)
				break
			}
		}
	}
}

func TestDeleteFolder(t *testing.T) {
	setupMockData()
	RegisterUser("user1")
	CreateFolder("user1", "folder1", "desc1")

	tests := []struct {
		username   string
		foldername string
		expected   error
	}{
		{"user1", "folder1", nil},
		{"user1", "folder1", errorDoesntExisted("folder1")},
		{"user2", "folder1", errorDoesntExisted("user2")},
	}

	for _, test := range tests {
		err := DeleteFolder(test.username, test.foldername)
		if err != nil && err.Error() != test.expected.Error() {
			t.Errorf("DeleteFolder(%s, %s) = %v; expected %v", test.username, test.foldername, err, test.expected)
		}
	}
}

func TestDeleteFile(t *testing.T) {
	setupMockData()
	RegisterUser("user1")
	CreateFolder("user1", "folder1", "desc1")
	CreateFile("user1", "folder1", "file1", "desc1")

	tests := []struct {
		username   string
		foldername string
		filename   string
		expected   error
	}{
		{"user1", "folder1", "file1", nil},
		{"user1", "folder1", "file1", errorDoesntExisted("file1")},
		{"user1", "folder2", "file1", errorDoesntExisted("folder2")},
		{"user2", "folder1", "file1", errorDoesntExisted("user2")},
	}

	for _, test := range tests {
		err := DeleteFile(test.username, test.foldername, test.filename)
		if err != nil && err.Error() != test.expected.Error() {
			t.Errorf("DeleteFile(%s, %s, %s) = %v; expected %v", test.username, test.foldername, test.filename, err, test.expected)
		}
	}
}

func TestRenameFolder(t *testing.T) {
	setupMockData()
	RegisterUser("user1")
	CreateFolder("user1", "folder1", "desc1")

	tests := []struct {
		username      string
		foldername    string
		newFolderName string
		expected      error
	}{
		{"user1", "folder1", "folder2", nil},
		{"user1", "folder1", "folder2", errorDoesntExisted("folder1")},
		{"user1", "folder2", "folder1", errorAlreayExisted("folder1")},
		{"user2", "folder1", "folder2", errorDoesntExisted("user2")},
		{"user1", "folder1", "invalid/folder", errorInvalidChars("invalid/folder")},
	}

	for _, test := range tests {
		err := RenameFolder(test.username, test.foldername, test.newFolderName)
		if err != nil && err.Error() != test.expected.Error() {
			t.Errorf("RenameFolder(%s, %s, %s) = %v; expected %v", test.username, test.foldername, test.newFolderName, err, test.expected)
		}
	}
}
