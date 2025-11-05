package goonobjects

import (
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Themaytrix/goon/utils"
)

type Blob struct{}

func ReadObject(path string) {


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

	// find goon directory and add retun relative .goon dir
	currDir, _ := os.Getwd()
	goonPath, isgoon := utils.IsGoonRepo(currDir, ".goon")
	fmt.Println(isgoon)
	if isgoon {
		hash_str := fmt.Sprintf("%x", hash)
		// get the obj dir
		obj_dir := filepath.Join(goonPath, "objects",hash_str[:2])
		blob_file := filepath.Join(obj_dir, hash_str[2:])
		// create object directory
		err := os.MkdirAll(obj_dir,0755 )
    if err != nil{
      panic(err)
    }
		// create a buffer
		var b bytes.Buffer
		w := zlib.NewWriter(&b)
    w.Write(store)
    w.Close()
    fmt.Println(b.Bytes())

		// write to file
    err = os.WriteFile(blob_file,b.Bytes(),0644)
if err != nil{
      panic(err)
    }

	}
}
