# Virtual File System

This is a virtual file system implemented in Go, allowing users to create, list, rename, and delete folders and files.

## Features

- Register users
- Create folders and files
- List folders and files with optional sorting
- Rename folders
- Delete folders and files
- Input validation for usernames, folder names, and file names

## Usage

### Commands
0. **help**
   
   Shows all the commands.

1. **register [username]**

   Registers a new user with the specified username.
   ```sh
   register user
   ```
   ```sh
   register "user A"
   ```

2. **create-folder [username] [foldername] [description]**

   Creates a new folder for the specified user with an 
   optional description.
   ```sh
   create-folder user folderA
   ```
   ```sh
   create-folder user folderB "folderB description"
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

7. **delete-file [username] [foldername] [filename]**
   
    Deletes the specified file in the folder for the user.

8. **rename-folder [username] [foldername] [new-folder-name]**
    Renames the specified folder for the user.
