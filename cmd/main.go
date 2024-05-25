// main.go
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"virtual-file-system/internal"
)

// parseArgs parses the input command and splits it into arguments considering quotes
func parseArgs(input string) []string {
	var args []string
	var current string
	inQuotes := false
	for _, char := range input {
		switch char {
		case ' ':
			if inQuotes {
				current += string(char)
			} else if current != "" {
				args = append(args, current)
				current = ""
			}
		case '"':
			inQuotes = !inQuotes
		default:
			current += string(char)
		}
	}
	if current != "" {
		args = append(args, current)
	}
	return args
}

// quoteIfNeeded adds double quotes around a string if it contains spaces
func quoteIfNeeded(s string) string {
	if strings.Contains(s, " ") {
		return fmt.Sprintf("\"%s\"", s)
	}
	return s
}

func main() {
	if err := internal.LoadData(); err != nil {
		fmt.Println("Error loading data:", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Virtual File System REPL")
	fmt.Println("------------------------")
	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		args := parseArgs(input)
		if len(args) < 1 {
			fmt.Println("No command provided")
			continue
		}

		command := args[0]
		switch command {
		case "register":
			if len(args) != 2 {
				fmt.Println("Usage: register [username]")
				continue
			}
			username := args[1]
			err := internal.RegisterUser(username)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Add", quoteIfNeeded(username), "successfully.")
			}
		case "create-folder":
			if len(args) < 3 || len(args) > 4 {
				fmt.Println("Usage: create-folder [username] [foldername] [description]?")
				continue
			}
			username := args[1]
			foldername := args[2]
			description := ""
			if len(args) == 4 {
				description = args[3]
			}
			err := internal.CreateFolder(username, foldername, description)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Create folder", quoteIfNeeded(foldername), "successfully.")
			}
		case "create-file":
			if len(args) < 4 || len(args) > 5 {
				fmt.Println("Usage: create-file [username] [foldername] [filename] [description]?")
				continue
			}
			username := args[1]
			foldername := args[2]
			filename := args[3]
			description := ""
			if len(args) == 5 {
				description = args[4]
			}
			err := internal.CreateFile(username, foldername, filename, description)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Create file", quoteIfNeeded(filename), "successfully.")
			}
		case "list-folders":
			if len(args) != 2 && len(args) != 4 {
				fmt.Println("Usage: list-folders [username] [--sort-name|--sort-created] [asc|desc]")
				continue
			}
			username := args[1]
			sortBy := "name"
			order := "asc"
			if len(args) == 4 {
				sortBy = strings.TrimPrefix(args[2], "--sort-")
				order = args[3]
				if order != "asc" && order != "desc" {
					fmt.Println("Usage: list-folders [username] [--sort-name|--sort-created] [asc|desc]")
					continue
				}
			}
			folders, err := internal.ListFolders(username, sortBy, order)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			for _, folder := range folders {
				if folder.Description != "" {
					fmt.Printf("%s %s %s %s\n", quoteIfNeeded(folder.Name), quoteIfNeeded(folder.Description), folder.CreatedAt.Format("2006-01-02 15:04:05"), quoteIfNeeded(username))
				} else {
					fmt.Printf("%s %s %s\n", quoteIfNeeded(folder.Name), folder.CreatedAt.Format("2006-01-02 15:04:05"), quoteIfNeeded(username))
				}
			}
		case "list-files":
			if len(args) != 3 && len(args) != 5 {
				fmt.Println("Usage: list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]")
				continue
			}
			username := args[1]
			foldername := args[2]
			sortBy := "name"
			order := "asc"
			if len(args) == 5 {
				sortBy = strings.TrimPrefix(args[3], "--sort-")
				order = args[4]
				if order != "asc" && order != "desc" {
					fmt.Println("Usage: list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]")
					continue
				}
			}
			files, err := internal.ListFiles(username, foldername, sortBy, order)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			for _, file := range files {
				if file.Description != "" {
					fmt.Printf("%s %s %s %s\n", quoteIfNeeded(file.Name), quoteIfNeeded(file.Description), file.CreatedAt.Format("2006-01-02 15:04:05"), quoteIfNeeded(username))
				} else {
					fmt.Printf("%s %s %s\n", quoteIfNeeded(file.Name), file.CreatedAt.Format("2006-01-02 15:04:05"), quoteIfNeeded(username))
				}
			}
		case "delete-folder":
			if len(args) != 3 {
				fmt.Println("Usage: delete-folder [username] [foldername]")
				continue
			}
			username := args[1]
			foldername := args[2]
			err := internal.DeleteFolder(username, foldername)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Delete folder", quoteIfNeeded(foldername), "successfully.")
			}
		case "delete-file":
			if len(args) != 4 {
				fmt.Println("Usage: delete-file [username] [foldername] [filename]")
				continue
			}
			username := args[1]
			foldername := args[2]
			filename := args[3]
			err := internal.DeleteFile(username, foldername, filename)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Delete file", quoteIfNeeded(filename), "successfully.")
			}
		case "rename-folder":
			if len(args) != 4 {
				fmt.Println("Usage: rename-folder [username] [foldername] [new-folder-name]")
				continue
			}
			username := args[1]
			foldername := args[2]
			newFolderName := args[3]
			err := internal.RenameFolder(username, foldername, newFolderName)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Rename folder", quoteIfNeeded(foldername), "to", quoteIfNeeded(newFolderName), "successfully.")
			}
		case "exit":
			fmt.Println("Exiting REPL...")
			return
		default:
			fmt.Println("Unrecognized command")
		}
	}
}
