package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"github.com/otiai10/copy"
)

// getBackupPath gets path to store backups
func getAllBackupsPath() string{
	backupHeadPath := os.Getenv("GOPATH") // Init to GOPATH
	dirSlice := []string{"src", "github.com", "AndrewYinLi", "QuitHub", "backup"} // Path to QuitHub src code
	for _,dir := range dirSlice{
		backupHeadPath = path.Join(backupHeadPath, dir)
	}
	//return "C:/Users/andre/Desktop/GoPath/backup" // for debugging
	return backupHeadPath
}

// Commit to the history for cd a copy of the cwd renamed as commitName
func commit(cd, commitName string){
	// Get paths
	backupHeadPath := path.Join(getAllBackupsPath(), filepath.Base(cd))
	backupCommitPath := path.Join(backupHeadPath, commitName)
	deleteCommit(cd, commitName) // If backupCommitPath exists, delete it and its contents
	copyToNewDir(backupCommitPath, cd, backupCommitPath) // Create backupCommitPath and copy all files from cd into backupCommitPath
}

// Creates empty directory at newPath and then copies contents of src to dest
func copyToNewDir(newPath, src, dest string){
	// Create newPath
	err := os.MkdirAll(newPath, os.ModePerm)
	if err != nil{
		log.Fatal(err)
	}
	// Copy all files from src to dest
	err = copy.Copy(src, dest)
	if err != nil{
		log.Fatal(err)
	}
}

// Revert the contents of the cwd to the contents of commitName stored in the history for cd
func revertCommit(cd, commitName string){
	// Get paths
	backupHeadPath := path.Join(getAllBackupsPath(), filepath.Base(cd))
	backupCommitPath := path.Join(backupHeadPath, commitName)
	// Check if backupCommitPath exists
	_, err := os.Stat(backupCommitPath);
	if os.IsNotExist(err) {
		log.Fatal("Commit to revert to does not exist.")
	}
	// Delete contents of cd and copy all files from backupCommitPath to cd
	deleteDir(cd)
	copyToNewDir(cd, backupCommitPath, cd)

}

// Print the history of commits for cd (which is really just the contents of the dir lol)
func commitHistory(cd string){
	backupHeadPath := path.Join(getAllBackupsPath(), filepath.Base(cd))
	// Open backupHeadPath
	f, err := os.Open(backupHeadPath)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	// Get and print file descriptions for all files in backupHeadPath
	files, err := f.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		fmt.Println(file.Name() + " " + file.ModTime().String())
	}
}

// Delete the directory commitName committed to cd
func deleteCommit(cd, commitName string){
	// Get paths
	backupHeadPath := path.Join(getAllBackupsPath(), filepath.Base(cd))
	backupCommitPath := path.Join(backupHeadPath, commitName)
	// If backupCommitPath exists, delete it and its contents
	deleteDir(backupCommitPath)
}

// Deletes directory and directory's contents
func deleteDir(dirPath string){
	_, err := os.Stat(dirPath);
	if !os.IsNotExist(err) {
		err := os.RemoveAll(dirPath)
		if err != nil{
			log.Fatal(err)
		}
	}
}

func main() {
	// Get args
	cd,_ := os.Getwd()
	//baseName := filepath.Base(cd)
	commitName := filepath.Base(cd)
	if len(os.Args) == 3{
		commitName = os.Args[2]
	}
	// Determine action
	if os.Args[1] == "commit"{
		commit(cd, commitName)
	} else if os.Args[1] == "revert"{
		revertCommit(cd, commitName)
	} else if os.Args[1] == "history"{
		commitHistory(cd)
	} else if os.Args[1] == "delete"{
		deleteCommit(cd, commitName)
	}
}
