package utils

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
)

func IsGoonRepo(currDir, targetDir string) (string,bool) {
	// check if .goon already exists in any path of the existing director

	if currDir == "/" {
		return "",false
	}

	targetPath := filepath.Join(currDir, targetDir)

	//  check if this directory exists
	_, err := os.Stat(targetPath)
	if err == nil {

		return targetPath,true
	}

	parentDir := filepath.Dir(currDir)
	return IsGoonRepo(parentDir, targetDir)
}

func HashBlob(path string) ([32]byte, error){
  var blob [32]byte

  // read the contents in the file
  content, err := os.ReadFile(path)

  if err != nil{
    return blob , err
  }

// define the header

  header := fmt.Sprintf("blob %d\x00", len(content))

  hash := sha1.New()
  hash.Write([]byte(header))
  hash.Write(content)

  copy(blob[:], hash.Sum(nil))

  return blob, nil
}


func Extractfiles(){

}

func IsFile(path string) bool {

  info, err := os.Stat(path)
  if err != nil{
    return false
  }
   // check if regular file

  return info.Mode().IsRegular()
}

func IsDir(path string) bool {
  info, err := os.Stat(path)
  if err != nil{
    return false
  }

  return info.IsDir()
}

func WalkDir(files *[]string, root string) error {
  // traverse root

  return filepath.WalkDir(root,func(path string, d os.DirEntry, err error) error {
    if err != nil{
      return err
    }

    // check if it's directory and not .goon or .git
    if d.IsDir() && d.Name() == ".git" && d.Name() == ".goon"{
      return filepath.SkipDir
    }

    if !d.IsDir(){
      *files = append(*files,path)
    }

    return nil

  })
}
