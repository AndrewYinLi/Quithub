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
	return "C:/Users/andre/Desktop/GoPath/backup" // for debugging
	//return backupHeadPath
}

// Commit to the history for cd a copy of the cwd renamed as commitName
func commit(cd, commitName string){
	// Get paths
	backupHeadPath := path.Join(getAllBackupsPath(), filepath.Base(cd))
	backupCommitPath := path.Join(backupHeadPath, commitName)
	// If backupCommitPath exists, delete it and its contents
	delete(cd, commitName)
	// Create backupCommitPath
	err := os.MkdirAll(backupCommitPath, os.ModePerm)
	if err != nil{
		log.Fatal(err)
	}
	// Copy all files from cd into backupCommitPath
	err = copy.Copy(cd, backupCommitPath)
	if err != nil{
		log.Fatal(err)
	}
}

// Revert the contents of the cwd to the contents of commitName stored in the history for cd
func revert(cd, commitName string){

}

// Print the history of commits for cd (which is really just the contents of the dir lol)
func history(cd string){
	backupHeadPath := path.Join(getAllBackupsPath(), filepath.Base(cd))
	f, err := os.Open(backupHeadPath)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	files, err := f.Readdir(-1)

	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		fmt.Println(file.Name() + " " + file.ModTime().String())
	}
}

// Delete the directory commitName committed to cd
func delete(cd, commitName string){
	// Get paths
	backupHeadPath := path.Join(getAllBackupsPath(), filepath.Base(cd))
	backupCommitPath := path.Join(backupHeadPath, commitName)
	// If backupCommitPath exists, delete it and its contents
	_, err := os.Stat(backupCommitPath);
	if !os.IsNotExist(err) {
		err := os.RemoveAll(backupCommitPath)
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
		revert(cd, commitName)
	} else if os.Args[1] == "history"{
		history(cd)
	} else if os.Args[1] == "delete"{
		delete(cd, commitName)
	}
}
