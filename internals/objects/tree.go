package tree

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Themaytrix/goon/internals/index"
)

type TreeEntry struct {
	Mode string
	Name string
	Hash [20]byte
}

type Tree struct {
	Entries []TreeEntry
}

type StackEntry struct {
	Name string
	Tree *Tree
}

// function builds tree from index files
func WriteTreefromIndex(entries []index.IndexEntry) ([20]byte,error) {
	root := &Tree{}        // intialize tree struct
	stack := []*Tree{root} // create a stack of Trees pointing to same tree and initalize first entry as address to root. fyi stack tracks our directories
	var prevparts []string

	for _, e := range entries {
		parts := strings.Split(e.Path, "/")
// check if directory already exist in our stack or directory is common
    common := 0
    for common <len(prevparts) && common < len(parts) && prevparts[common] == parts[common]{
      common ++
    }

    // close the directories we already visisted
    for len(stack)-1 > common{
      tree := stack[len(stack)-1]
      // pop top element
      stack = stack[:len(stack)-1]
      // write the tree object
      hash,err := WriteTreeObject(tree.Entries)
      if err != nil{
        return [20]byte{},err
    }
      //track parent folder for file eg. src/utils/file.txt. we get the parent utils
      parent := stack[len(stack)-1]
      
      // populating a tree entry
      parent.Entries = append(parent.Entries, TreeEntry{
        Mode: "40000",
        Name: prevparts[len(stack)-1],
        Hash: hash,
      })


    }
  
    // handling new directories like src/main.go
    
    for i:= common; i<len(parts)-1; i++{
    // create a new tree and add to our stack-remember stack tracks our directories  
      t := &Tree{}
      stack = append(stack,t)
    }

// add entry. get the current dir
    current := stack[len(stack)-1]

    current.Entries = append(current.Entries, TreeEntry{
      Mode: "100644",
      Name: parts[len(parts)-1],
      Hash: e.ObjectID,
    })

prevparts = parts
  }

    for len(stack) >1{
      tree := stack[len(stack)-1]
      stack = stack[:len(stack)-1]
    hash,err := WriteTreeObject(tree.Entries)
    if err!=nil{
      return [20]byte{},err
    }
    parent := stack[len(stack)-1]

    parent.Entries = append(parent.Entries, TreeEntry{
      Mode: "40000",
      Name: prevparts[len(stack)-1],
      Hash: hash,
    })

    }
return WriteTreeObject(root.Entries)		
}

func WriteTreeObject(entries []TreeEntry)([20]byte, error){
var b bytes.Buffer
  for _,e := range entries{
    fmt.Fprintf(&b,"%s %s",e.Mode, e.Name)
    b.WriteByte(0)
    b.Write(e.Hash[:])
  }
  return WriteObject("tree",b.Bytes())
}


func WriteObject(objType string, data []byte)([20]byte,error){

}
