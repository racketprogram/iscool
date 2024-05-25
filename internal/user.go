// internal/user.go
package internal

import (
	"fmt"
	"sort"
	"time"
)

// RegisterUser registers a new user with a unique username
func RegisterUser(username string) error {
	if _, exists := users[username]; exists {
		return fmt.Errorf("The %s has already existed.", username)
	}
	if !isValidUsername(username) {
		return fmt.Errorf("The %s contains invalid chars.", username)
	}
	users[username] = &User{
		Username: username,
		Folders:  make(map[string]*Folder),
	}
	return SaveData()
}

// CreateFolder creates a new folder for a user
func CreateFolder(username, foldername string, description ...string) error {
	user, exists := users[username]
	if !exists {
		return fmt.Errorf("The %s doesn't exist.", username)
	}

	if _, exists := user.Folders[foldername]; exists {
		return fmt.Errorf("The %s has already existed.", foldername)
	}

	desc := ""
	if len(description) > 0 {
		desc = description[0]
	}

	user.Folders[foldername] = &Folder{
		Name:        foldername,
		Description: desc,
		CreatedAt:   time.Now(),
		Files:       make(map[string]*File),
	}
	return SaveData()
}

// CreateFile creates a new file in a user's folder
func CreateFile(username, foldername, filename string, description ...string) error {
	user, exists := users[username]
	if !exists {
		return fmt.Errorf("The %s doesn't exist.", username)
	}

	folder, exists := user.Folders[foldername]
	if !exists {
		return fmt.Errorf("The %s doesn't exist.", foldername)
	}

	if !isValidFilename(filename) {
		return fmt.Errorf("The %s contains invalid chars.", filename)
	}
	if _, exists := folder.Files[filename]; exists {
		return fmt.Errorf("The %s has already existed.", filename)
	}

	desc := ""
	if len(description) > 0 {
		desc = description[0]
	}

	folder.Files[filename] = &File{
		Name:        filename,
		Description: desc,
		CreatedAt:   time.Now(),
	}
	return SaveData()
}

// ListFolders lists all folders for a user with optional sorting
func ListFolders(username, sortBy, order string) ([]*Folder, error) {
	user, exists := users[username]
	if !exists {
		return nil, fmt.Errorf("The %s doesn't exist.", username)
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
				return folders[i].CreatedAt.After(folders[j].CreatedAt)
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
		return nil, fmt.Errorf("The %s doesn't exist.", username)
	}

	folder, exists := user.Folders[foldername]
	if !exists {
		return nil, fmt.Errorf("The %s doesn't exist.", foldername)
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
				return files[i].CreatedAt.After(files[j].CreatedAt)
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
		return fmt.Errorf("The %s doesn't exist.", username)
	}

	if _, exists := user.Folders[foldername]; !exists {
		return fmt.Errorf("The %s doesn't exist.", foldername)
	}

	delete(user.Folders, foldername)
	return SaveData()
}

// DeleteFile deletes a file in a user's folder
func DeleteFile(username, foldername, filename string) error {
	user, exists := users[username]
	if !exists {
		return fmt.Errorf("The %s doesn't exist.", username)
	}

	folder, exists := user.Folders[foldername]
	if !exists {
		return fmt.Errorf("The %s doesn't exist.", foldername)
	}

	if _, exists := folder.Files[filename]; !exists {
		return fmt.Errorf("The %s doesn't exist.", filename)
	}

	delete(folder.Files, filename)
	return SaveData()
}

// RenameFolder renames a folder for a user
func RenameFolder(username, foldername, newFolderName string) error {
	user, exists := users[username]
	if !exists {
		return fmt.Errorf("The %s doesn't exist.", username)
	}

	if _, exists := user.Folders[foldername]; !exists {
		return fmt.Errorf("The %s doesn't exist.", foldername)
	}
	if _, exists := user.Folders[newFolderName]; exists {
		return fmt.Errorf("The %s has already existed.", newFolderName)
	}

	folder := user.Folders[foldername]
	folder.Name = newFolderName
	user.Folders[newFolderName] = folder
	delete(user.Folders, foldername)
	return SaveData()
}

// isValidUsername validates the username
func isValidUsername(username string) bool {
	// Add username validation logic
	return true
}

// isValidFilename validates the filename
func isValidFilename(filename string) bool {
	// Add filename validation logic
	return true
}
