// internal/types.go
package internal

import "time"

// User represents a system user
type User struct {
	Username string             `json:"username"`
	Folders  map[string]*Folder `json:"folders"`
}

// Folder represents a folder in the file system
type Folder struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	CreatedAt   time.Time        `json:"created_at"`
	Files       map[string]*File `json:"files"`
}

// File represents a file in the file system
type File struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
