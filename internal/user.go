// internal/user.go
package internal

import (
	"fmt"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"
)

// windowsSleep function
//
// This function introduces a nanosecond-level sleep when running on Windows systems.
// The primary purpose of this function is to address the low time resolution issue on Windows.
// On Windows, the system time is typically updated every 10-15 milliseconds. Therefore, querying the current time multiple times within this period might return the same value.
// This can cause issues in scenarios requiring high-precision timestamps, such as in tests.
//
// By introducing a brief sleep, we ensure that consecutive calls to get the current time will yield different values, thereby providing unique timestamps.
// This function is a no-op on non-Windows systems.
func windowsSleep() {
	if runtime.GOOS == "windows" {
		time.Sleep(time.Nanosecond)
	}
}

// QuoteIfNeeded adds double quotes around a string if it contains spaces
func QuoteIfNeeded(s string) string {
	if strings.Contains(s, " ") {
		return fmt.Sprintf("\"%s\"", s)
	}
	return s
}

// isValidName validates the name, allowing letters, numbers, spaces, underscores, and hyphens, with a length of 1-50 characters.
func isValidName(name string) bool {
	// Define the regular expression for a valid name
	var validNameRegex = regexp.MustCompile(`^[a-zA-Z0-9 _-]{1,50}$`)
	return validNameRegex.MatchString(name)
}

func errorAlreayExisted(name string) error {
	return fmt.Errorf("The %s has already existed.", QuoteIfNeeded(name))
}

func errorDoesntExisted(name string) error {
	return fmt.Errorf("The %s doesn't exist.", QuoteIfNeeded(name))
}

func errorInvalidChars(name string) error {
	return fmt.Errorf("The %s contains invalid chars.", QuoteIfNeeded(name))
}

// RegisterUser registers a new user with a unique username
func RegisterUser(username string) error {
	if _, exists := users[username]; exists {
		return errorAlreayExisted(username)
	}
	if !isValidName(username) {
		return errorInvalidChars(username)
	}
	users[username] = &User{
		Username: username,
		Folders:  make(map[string]*Folder),
	}
	return SaveData()
}

// CreateFolder creates a new folder for a user
func CreateFolder(username, foldername string, description string) error {
	user, exists := users[username]
	if !exists {
		return errorDoesntExisted(username)
	}

	if _, exists := user.Folders[foldername]; exists {
		return errorAlreayExisted(foldername)
	}

	if !isValidName(foldername) {
		return errorInvalidChars(foldername)
	}

	user.Folders[foldername] = &Folder{
		Name:        foldername,
		Description: description,
		CreatedAt:   time.Now(),
		Files:       make(map[string]*File),
	}
	windowsSleep()
	return SaveData()
}

// CreateFile creates a new file in a user's folder
func CreateFile(username, foldername, filename string, description string) error {
	user, exists := users[username]
	if !exists {
		return errorDoesntExisted(username)
	}

	folder, exists := user.Folders[foldername]
	if !exists {
		return errorDoesntExisted(foldername)
	}

	if !isValidName(filename) {
		return errorInvalidChars(filename)
	}

	if _, exists := folder.Files[filename]; exists {
		return errorAlreayExisted(filename)
	}

	folder.Files[filename] = &File{
		Name:        filename,
		Description: description,
		CreatedAt:   time.Now(),
	}
	windowsSleep()
	return SaveData()
}

// ListFolders lists all folders for a user with optional sorting
func ListFolders(username, sortBy, order string) ([]*Folder, error) {
	user, exists := users[username]
	if !exists {
		return nil, errorDoesntExisted(username)
	}

	folders := make([]*Folder, 0, len(user.Folders))
	for _, folder := range user.Folders {
		folders = append(folders, folder)
	}

	// Default sorting by name in ascending order
	if sortBy == "" {
		sortBy = "name"
	}
	if order == "" {
		order = "asc"
	}

	switch sortBy {
	case "name":
		sort.Slice(folders, func(i, j int) bool {
			if order == "desc" {
				return folders[i].Name > folders[j].Name
			}
			return folders[i].Name < folders[j].Name
		})
	case "created":
		sort.Slice(folders, func(i, j int) bool {
			if order == "desc" {
				return !folders[i].CreatedAt.Before(folders[j].CreatedAt)
			}
			return folders[i].CreatedAt.Before(folders[j].CreatedAt)
		})
	default:
		sort.Slice(folders, func(i, j int) bool {
			return folders[i].Name < folders[j].Name
		})
	}

	return folders, nil
}

// ListFiles lists all files in a user's folder with optional sorting
func ListFiles(username, foldername, sortBy, order string) ([]*File, error) {
	user, exists := users[username]
	if !exists {
		return nil, errorDoesntExisted(username)
	}

	folder, exists := user.Folders[foldername]
	if !exists {
		return nil, errorDoesntExisted(foldername)
	}

	files := make([]*File, 0, len(folder.Files))
	for _, file := range folder.Files {
		files = append(files, file)
	}

	// Default sorting by name in ascending order
	if sortBy == "" {
		sortBy = "name"
	}
	if order == "" {
		order = "asc"
	}

	switch sortBy {
	case "name":
		sort.Slice(files, func(i, j int) bool {
			if order == "desc" {
				return files[i].Name > files[j].Name
			}
			return files[i].Name < files[j].Name
		})
	case "created":
		sort.Slice(files, func(i, j int) bool {
			if order == "desc" {
				return !files[i].CreatedAt.Before(files[j].CreatedAt)
			}
			return files[i].CreatedAt.Before(files[j].CreatedAt)
		})
	default:
		sort.Slice(files, func(i, j int) bool {
			return files[i].Name < files[j].Name
		})
	}

	return files, nil
}

// DeleteFolder deletes a folder for a user
func DeleteFolder(username, foldername string) error {
	user, exists := users[username]
	if !exists {
		return errorDoesntExisted(username)
	}

	if _, exists := user.Folders[foldername]; !exists {
		return errorDoesntExisted(foldername)
	}

	delete(user.Folders, foldername)
	return SaveData()
}

// DeleteFile deletes a file in a user's folder
func DeleteFile(username, foldername, filename string) error {
	user, exists := users[username]
	if !exists {
		return errorDoesntExisted(username)
	}

	folder, exists := user.Folders[foldername]
	if !exists {
		return errorDoesntExisted(foldername)
	}

	if _, exists := folder.Files[filename]; !exists {
		return errorDoesntExisted(filename)
	}

	delete(folder.Files, filename)
	return SaveData()
}

// RenameFolder renames a folder for a user
func RenameFolder(username, foldername, newFolderName string) error {
	user, exists := users[username]
	if !exists {
		return errorDoesntExisted(username)
	}

	if _, exists := user.Folders[foldername]; !exists {
		return errorDoesntExisted(foldername)
	}

	if !isValidName(newFolderName) {
		return errorInvalidChars(newFolderName)
	}

	if _, exists := user.Folders[newFolderName]; exists {
		return errorAlreayExisted(newFolderName)
	}

	folder := user.Folders[foldername]
	folder.Name = newFolderName
	user.Folders[newFolderName] = folder
	delete(user.Folders, foldername)
	return SaveData()
}
