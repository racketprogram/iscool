# Virtual File System

This is a virtual file system implemented in Go, allowing users to create, list, rename, and delete folders and files. The system uses JSON for data storage, ensuring a simple and portable way to manage the file system's state. The JSON file, data.json, is used to persist user data, folders, and files across sessions.

## Design Principles
- Simplicity: Uses JSON for data storage, providing a straightforward and human-readable format.
- Portability: The system can be easily moved or shared, as it relies on a single JSON file for data persistence.
- Robustness: Designed to handle typical file system operations with clear input validation to ensure data integrity.
- Ease of Use: Provides a simple REPL (Read-Eval-Print Loop) interface for interacting with the virtual file system.

## Features

- Help
- Register users
- Create folders and files
- List folders and files with optional sorting
- Rename folders
- Delete folders and files
- Input validation for usernames, folder names, and file names

## Build

To build the project, you need to have Go installed on your machine. Follow the instructions below to clone the repository and build the executable.

1. Clone the repository:
   ```sh
   git clone https://github.com/racketprogram/iscool.git
   cd iscool
   ```

2. Build the project:
    ```sh
    go build -o vfs ./cmd/main.go
    ```
3. Run the executable:
    ```sh
    ./vfs
    ```

## Test
### Integration test
```sh
go test ./cmd/...
```
### Unit test
```sh
go test ./internal/...
```

## Usage

The virtual file system is implemented as a REPL (Read-Eval-Print Loop) where users can type commands to interact with the system. Usernames, folder names, and file names are case-insensitive. Below are the available commands and their usage.

### Specifying Data File Path

You can specify the location and name of the `data.json` file by passing the path as a command-line argument when starting the application. If the provided path is a directory or ends with a `/`, an error will be returned. By default, the `data.json` file will be created and used in the directory where the application is executed.

#### Example:

```sh
./vfs /path/to/custom_data.json
./vfs /path/to/directory/ ❌ # This will return an error
./vfs /path/to/directory ❌ # This will return an error
```

In the first example, the data will be stored in /path/to/custom_data.json. In the second and third examples, an error will be returned since the provided path is a directory or ends with a /.

### Commands
0. **help**
   
   Shows all the commands.
   ```sh
   > help
    Usage: register [username]
    Usage: create-folder [username] [foldername] [description]?
    Usage: create-file [username] [foldername] [filename] [description]?
    Usage: list-folders [username] [--sort-name|--sort-created] [asc|desc]
    Usage: list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]
    Usage: delete-folder [username] [foldername]
    Usage: delete-file [username] [foldername] [filename]
    Usage: rename-folder [username] [foldername] [new-folder-name]

    [username] [foldername] and [filename] are case insensitive.
   ```

1. **register [username]**

   Registers a new user with the specified username.
   ```sh
   register userA # It will actually be stored as "usera".
   ```
   ```sh
   register "user A" # It will actually be stored as "user a".
   ```

2. **create-folder [username] [foldername] [description]**

   Creates a new folder for the specified user with an 
   optional description.
   ```sh
   create-folder user folderA
   ```
   ```sh
   create-folder user folderA "folderB description"
   ```
   ```sh
   create-folder "user A" "folder A" "folder A description"
   ```
3. **create-file [username] [foldername] [filename] [description]**

    Creates a new file in the specified folder for the user with an optional description.
    ```sh
    create-file user folderA fileA
    ```
    ```sh
    create-file user folderA fileB "fileB description"
    ```
    ```sh
    create-file "user A" "folder A" "file A" "file A description"
    ```

4. **list-folders [username] [--sort-name|--sort-created] [asc|desc]**

    Lists all folders for the specified user with optional sorting.

    ```sh
    list-folders user
    ```
    ```sh
    list-folders user --sort-name asc
    ```
    ```sh
    list-folders user --sort-created desc
    ```
    ```sh
    list-folders user --sort-created ❌ # The order is necessary when specifying sort criteria.
    ```


5. **list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]**

    Lists all folders for the specified user with optional sorting.

    ```sh
    list-files user folderA
    ```
    ```sh
    list-files user folderA --sort-name asc
    ```
    ```sh
    list-files user folderA --sort-created desc
    ```
    ```sh
    list-files user folderA --sort-created ❌ # The order is necessary when specifying sort criteria.
    ```

6. **delete-folder [username] [foldername]**

    Deletes the specified folder for the user.

    ```sh
    delete-folder user folderA
    ```
    ```sh
    delete-folder "user A" "folder A"
    ```

7. **delete-file [username] [foldername] [filename]**
   
    Deletes the specified file in the folder for the user.

    ```sh
    delete-file user folderA fileA
    ```
    ```sh
    delete-file "user A" "folder A" "file A"
    ```

8. **rename-folder [username] [foldername] [new-folder-name]**
   
    Renames the specified folder for the user.

    ```sh
    rename-folder user folderA newFolderName
    ```
    ```sh
    rename-folder "user A" "folder A" "new folder name"
    ```

## Input Validation Rules

### Usernames:

- Can contain letters, numbers, spaces, underscores (`_`), and hyphens (`-`).
- Length: 1-50 characters.

### Folder names:

- Can contain letters, numbers, spaces, underscores (`_`), and hyphens (`-`).
- Length: 1-50 characters.

### File names:

- Can contain letters, numbers, spaces, underscores (`_`), and hyphens (`-`).
- Length: 1-50 characters.
