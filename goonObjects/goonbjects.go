package goonobjects

import (
	"os"
  "fmt"
)

type Blob struct{}

func readObject() {
}

func hashObject(file string) {
	// read contents
  content, err := os.ReadFile(file)
  if err != nil {
    panic(err)
  }
  fmt.Printf("%T",content)
	// add the type
	// hash contents
}
