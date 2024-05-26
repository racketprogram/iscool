// main.go
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"virtual-file-system/internal"
)

var quoteIfNeeded = internal.QuoteIfNeeded

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

// handleCommand processes a single command
func handleCommand(command string, args []string) {
	switch command {
	case "register":
		if len(args) != 1 {
			fmt.Println("Usage: register [username]")
			return
		}
		username := args[0]
		err := internal.RegisterUser(username)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Add", quoteIfNeeded(username), "successfully.")
		}
	case "create-folder":
		if len(args) < 2 || len(args) > 3 {
			fmt.Println("Usage: create-folder [username] [foldername] [description]?")
			return
		}
		username := args[0]
		foldername := args[1]
		description := ""
		if len(args) == 3 {
			description = args[2]
		}
		err := internal.CreateFolder(username, foldername, description)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Create", quoteIfNeeded(foldername), "successfully.")
		}
	case "create-file":
		if len(args) < 3 || len(args) > 4 {
			fmt.Println("Usage: create-file [username] [foldername] [filename] [description]?")
			return
		}
		username := args[0]
		foldername := args[1]
		filename := args[2]
		description := ""
		if len(args) == 4 {
			description = args[3]
		}
		err := internal.CreateFile(username, foldername, filename, description)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Create", quoteIfNeeded(filename), "successfully.")
		}
	case "list-folders":
		if len(args) != 1 && len(args) != 3 {
			fmt.Println("Usage: list-folders [username] [--sort-name|--sort-created] [asc|desc]")
			return
		}
		username := args[0]
		sortBy := "name"
		order := "asc"
		if len(args) == 3 {
			sortBy = strings.TrimPrefix(args[1], "--sort-")
			order = args[2]
			if order != "asc" && order != "desc" {
				fmt.Println("Usage: list-folders [username] [--sort-name|--sort-created] [asc|desc]")
				return
			}
		}
		folders, err := internal.ListFolders(username, sortBy, order)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		for _, folder := range folders {
			if folder.Description != "" {
				fmt.Printf("%s %s %s %s\n", quoteIfNeeded(folder.Name), quoteIfNeeded(folder.Description), folder.CreatedAt.Format("2006-01-02 15:04:05"), quoteIfNeeded(username))
			} else {
				fmt.Printf("%s %s %s\n", quoteIfNeeded(folder.Name), folder.CreatedAt.Format("2006-01-02 15:04:05"), quoteIfNeeded(username))
			}
		}
	case "list-files":
		if len(args) != 2 && len(args) != 4 {
			fmt.Println("Usage: list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]")
			return
		}
		username := args[0]
		foldername := args[1]
		sortBy := "name"
		order := "asc"
		if len(args) == 4 {
			sortBy = strings.TrimPrefix(args[2], "--sort-")
			order = args[3]
			if order != "asc" && order != "desc" {
				fmt.Println("Usage: list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]")
				return
			}
		}
		files, err := internal.ListFiles(username, foldername, sortBy, order)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		for _, file := range files {
			if file.Description != "" {
				fmt.Printf("%s %s %s %s\n", quoteIfNeeded(file.Name), quoteIfNeeded(file.Description), file.CreatedAt.Format("2006-01-02 15:04:05"), quoteIfNeeded(username))
			} else {
				fmt.Printf("%s %s %s\n", quoteIfNeeded(file.Name), file.CreatedAt.Format("2006-01-02 15:04:05"), quoteIfNeeded(username))
			}
		}
	case "delete-folder":
		if len(args) != 2 {
			fmt.Println("Usage: delete-folder [username] [foldername]")
			return
		}
		username := args[0]
		foldername := args[1]
		err := internal.DeleteFolder(username, foldername)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Delete", quoteIfNeeded(foldername), "successfully.")
		}
	case "delete-file":
		if len(args) != 3 {
			fmt.Println("Usage: delete-file [username] [foldername] [filename]")
			return
		}
		username := args[0]
		foldername := args[1]
		filename := args[2]
		err := internal.DeleteFile(username, foldername, filename)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Delete", quoteIfNeeded(filename), "successfully.")
		}
	case "rename-folder":
		if len(args) != 3 {
			fmt.Println("Usage: rename-folder [username] [foldername] [new-folder-name]")
			return
		}
		username := args[0]
		foldername := args[1]
		newFolderName := args[2]
		err := internal.RenameFolder(username, foldername, newFolderName)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Rename", quoteIfNeeded(foldername), "to", quoteIfNeeded(newFolderName), "successfully.")
		}
	case "exit":
		fmt.Println("Exiting REPL...")
		os.Exit(0)
	default:
		fmt.Println("Unrecognized command")
	}
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
		handleCommand(command, args[1:])
	}
}
