package goonobjects

import (
	"crypto/sha256"
	"fmt"
	"os"
  "github.com/Themaytrix/goon/utils"
)

type Blob struct{}

func readObject() {
}

func HashObject(file string) {
	// read contents
	content, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", content)
	// define header
	header := fmt.Sprintf("blob %d\u0000", len(content))
	fmt.Println(header)
	store := append([]byte(header), content...)
	// hash contents
  hash := sha256.Sum256(store)

//find goon directory and add retun relative .goon dir
  currDir, _ := os.Getwd()
  goonPath,isgoon := utils.IsGoonRepo(currDir,".goon")
  fmt.Println(isgoon)
  if isgoon{
    fmt.Println(goonPath)
    fmt.Printf("%x \n",hash)
  }


}
