// main.go
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"virtual-file-system/internal"
)

var caseInsensitive = true

var quoteIfNeeded = internal.QuoteIfNeeded

var commnadRegister = "Usage: register [username]"
var commnadCreateFolder = "Usage: create-folder [username] [foldername] [description]?"
var commnadCreateFile = "Usage: create-file [username] [foldername] [filename] [description]?"
var commnadListFolders = "Usage: list-folders [username] [--sort-name|--sort-created] [asc|desc]"
var commnadListFiles = "Usage: list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]"
var commandDeleteFolder = "Usage: delete-folder [username] [foldername]"
var commandDeleteFile = "Usage: delete-file [username] [foldername] [filename]"
var commandRenameFolder = "Usage: rename-folder [username] [foldername] [new-folder-name]"
var commands = []string{
	commnadRegister,
	commnadCreateFolder,
	commnadCreateFile,
	commnadListFolders,
	commnadListFiles,
	commandDeleteFolder,
	commandDeleteFile,
	commandRenameFolder,
}

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
			fmt.Println(commnadRegister)
			return
		}
		username := args[0]
		if caseInsensitive {
			username = strings.ToLower(username)
		}
		err := internal.RegisterUser(username)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
		} else {
			fmt.Println("Add", quoteIfNeeded(username), "successfully.")
		}
	case "create-folder":
		if len(args) < 2 || len(args) > 3 {
			fmt.Println(commnadCreateFolder)
			return
		}
		username := args[0]
		foldername := args[1]
		if caseInsensitive {
			username = strings.ToLower(username)
			foldername = strings.ToLower(foldername)
		}
		description := ""
		if len(args) == 3 {
			description = args[2]
		}
		err := internal.CreateFolder(username, foldername, description)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
		} else {
			fmt.Println("Create", quoteIfNeeded(foldername), "successfully.")
		}
	case "create-file":
		if len(args) < 3 || len(args) > 4 {
			fmt.Println(commnadCreateFile)
			return
		}
		username := args[0]
		foldername := args[1]
		filename := args[2]
		if caseInsensitive {
			username = strings.ToLower(username)
			foldername = strings.ToLower(foldername)
			filename = strings.ToLower(filename)
		}
		description := ""
		if len(args) == 4 {
			description = args[3]
		}
		err := internal.CreateFile(username, foldername, filename, description)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
		} else {
			fmt.Printf("Create %s in %s/%s successfully.\n", quoteIfNeeded(filename), quoteIfNeeded(username), quoteIfNeeded(foldername))
		}
	case "list-folders":
		if len(args) != 1 && len(args) != 3 {
			fmt.Println(commnadListFolders)
			return
		}
		username := args[0]
		if caseInsensitive {
			username = strings.ToLower(username)
		}
		sortBy := "name"
		order := "asc"
		if len(args) == 3 {
			sortBy = strings.TrimPrefix(args[1], "--sort-")
			order = args[2]
			if order != "asc" && order != "desc" {
				fmt.Println(commnadListFolders)
				return
			}
		}
		folders, err := internal.ListFolders(username, sortBy, order)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			return
		}
		if len(folders) == 0 {
			fmt.Printf("Warning: The %s doesn't have any folders.\n", quoteIfNeeded(username))
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
			fmt.Println(commnadListFiles)
			return
		}
		username := args[0]
		foldername := args[1]
		if caseInsensitive {
			username = strings.ToLower(username)
			foldername = strings.ToLower(foldername)
		}
		sortBy := "name"
		order := "asc"
		if len(args) == 4 {
			sortBy = strings.TrimPrefix(args[2], "--sort-")
			order = args[3]
			if order != "asc" && order != "desc" {
				fmt.Println(commnadListFiles)
				return
			}
		}
		files, err := internal.ListFiles(username, foldername, sortBy, order)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			return
		}
		if len(files) == 0 {
			fmt.Println("Warning: The folder is empty.")
			return
		}
		for _, file := range files {
			if file.Description != "" {
				fmt.Printf("%s %s %s %s %s\n", quoteIfNeeded(file.Name), quoteIfNeeded(file.Description), file.CreatedAt.Format("2006-01-02 15:04:05"), quoteIfNeeded(foldername), quoteIfNeeded(username))
			} else {
				fmt.Printf("%s %s %s %s\n", quoteIfNeeded(file.Name), file.CreatedAt.Format("2006-01-02 15:04:05"), quoteIfNeeded(foldername), quoteIfNeeded(username))
			}
		}
	case "delete-folder":
		if len(args) != 2 {
			fmt.Println(commandDeleteFolder)
			return
		}
		username := args[0]
		foldername := args[1]
		if caseInsensitive {
			username = strings.ToLower(username)
			foldername = strings.ToLower(foldername)
		}
		err := internal.DeleteFolder(username, foldername)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
		} else {
			fmt.Println("Delete", quoteIfNeeded(foldername), "successfully.")
		}
	case "delete-file":
		if len(args) != 3 {
			fmt.Println(commandDeleteFile)
			return
		}
		username := args[0]
		foldername := args[1]
		filename := args[2]
		if caseInsensitive {
			username = strings.ToLower(username)
			foldername = strings.ToLower(foldername)
			filename = strings.ToLower(filename)
		}
		err := internal.DeleteFile(username, foldername, filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
		} else {
			fmt.Printf("Delete %s in %s/%s successfully.\n", quoteIfNeeded(filename), quoteIfNeeded(username), quoteIfNeeded(foldername))
		}
	case "rename-folder":
		if len(args) != 3 {
			fmt.Println(commandRenameFolder)
			return
		}
		username := args[0]
		foldername := args[1]
		newFolderName := args[2]
		if caseInsensitive {
			username = strings.ToLower(username)
			foldername = strings.ToLower(foldername)
			newFolderName = strings.ToLower(newFolderName)
		}
		err := internal.RenameFolder(username, foldername, newFolderName)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
		} else {
			fmt.Println("Rename", quoteIfNeeded(foldername), "to", quoteIfNeeded(newFolderName), "successfully.")
		}
	case "exit":
		fmt.Println("Exiting REPL...")
		os.Exit(0)
	case "help":
		for _, command := range commands {
			fmt.Println(command)
		}
	default:
		fmt.Println("Unrecognized command")
	}
}

func main() {
	if err := internal.LoadData(); err != nil {
		fmt.Fprintln(os.Stderr, "Error: loading data:", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Virtual File System REPL")
	fmt.Println("Type `help` to show the commands.")
	fmt.Println("------------------------")
	for {
		fmt.Print("# ")
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
