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
  Count uint32
  Entries []*IndexEntry
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

// this function is to construct a new instance of the Entry
func NewEntry(path string) (*IndexEntry, error){
  // get file information
  fileInfo, err := os.Stat(path)
  if err != nil{
    return nil,err
  }
  // open file
  file,err := os.Open(path)
  if err != nil{
    return nil,err
  }
  defer file.Close()
  // read file contents
  h := sha256.New()
  io.Copy(h,file) // stream the filecontents for hashing. discard filebytes and retain whats necessary for hashing

  // hash file contents. execute file hash.
  var oid [32]byte
  copy(oid[:],h.Sum(nil)) // deep copy of hash for ObjectID in IndexEntry

  // get low level fileinformation for IndexEntry
  stat := fileInfo.Sys().(*syscall.Stat_t)
  // write new entry 
  entry := &IndexEntry{
    CTimeSec: uint32(stat.Ctimespec.Sec),
    CTimeNSec: uint32(stat.Ctimespec.Nsec),
    MTimeSec: uint32(stat.Mtimespec.Sec),
    MTimeNSec: uint32(stat.Mtimespec.Nsec),
    Dev: uint32(stat.Dev),
    Ino: uint32(stat.Ino),
    Mode: uint32(stat.Mode),
    UID: uint32(stat.Uid),
    GID: uint32(stat.Gid),
    FileSize: uint32(fileInfo.Size()),
    ObjectID: oid,
    Flags: uint16(len(filepath.Base(path))),
    Path: path,
  }

  return entry, nil

}


// write the Index entry into goon/index file
