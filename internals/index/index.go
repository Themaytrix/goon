package index

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
)

type Index struct{

  Signature [4]byte
  Version uint32
  Entries []IndexEntry
}

type IndexEntry struct{
CTimeSec  uint32
	CTimeNSec uint32
	MTimeSec  uint32
	MTimeNSec uint32
	Dev       uint32
	Ino       uint32
	Mode      uint32
	UID       uint32
	GID       uint32
	FileSize  uint32
	ObjectID  [32]byte // SHA256 hash of file contents
	Flags     uint16
	Path      string
}
// goon Index entry
func NewIndex() *Index{
  return &Index{
    Signature: [4]byte{'D','I','R','C'},
    Version: 2,
    Entries: []IndexEntry{},
  }
}

// add index entries to the index
func (idx *Index) AddEntry(e IndexEntry){
  idx.Entries = append(idx.Entries, e)
}


// create a new indexentry
func NewIndexEntry(path string)(IndexEntry, error){
  // check if the file exists

  info,err := os.Stat(path)

  if err != nil{
    return IndexEntry{},err
  }

  raw := info.Sys()

}
