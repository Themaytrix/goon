package index

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"syscall"

	"github.com/Themaytrix/goon/utils"
)

type Index struct {
	Signature [4]byte
	Version   uint32
	Entries   []IndexEntry
}

type IndexEntry struct {
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
func NewIndex() *Index {
	return &Index{
		Signature: [4]byte{'D', 'I', 'R', 'C'},
		Version:   2,
		Entries:   []IndexEntry{},
	}
}

// add index entries to the index
func (idx *Index) AddEntry(e IndexEntry) {
	idx.Entries = append(idx.Entries, e)
}

// create a new indexentry
func NewIndexEntry(path string) (IndexEntry, error) {
	// check if the file exists

	info, err := os.Stat(path)
	if err != nil {
		return IndexEntry{}, err
	}

	raw := info.Sys()
	cSec, cNSec, mSec, mNSec, dev, ino, uid, gid := extractMetadata(raw)

	// hash object
	blob, err := utils.HashBlob(path)
	if err != nil {
		// return empty IndexEntry
		return IndexEntry{}, err
	}

	// handling flags
	flags := uint16(len(path))

	if flags > 0x0FFF {
		flags = 0x0FFF
	}

	return IndexEntry{
		CTimeSec:  cSec,
		CTimeNSec: cNSec,
		MTimeSec:  mSec,
		MTimeNSec: mNSec,
		Dev:       dev,
		Ino:       ino,
		Mode:      uint32(0o100644),
		UID:       uid,
		GID:       gid,
		FileSize:  uint32(info.Size()),
		ObjectID:  blob,
		Flags:     flags,
		Path:      filepath.ToSlash(path),
	}, nil
}

// write new index to disk
func (idx *Index) WriteIndex(path string) error {
  // create a buffer to write format into memory before writing to disk or file

  var buf bytes.Buffer


  buf.Write(idx.Signature[:])
  binary.Write(&buf, binary.BigEndian, idx.Version)
  binary.Write(&buf, binary.BigEndian, uint32(len(idx.Entries)))

  // sort the index entries in lexographical order

  sort.Slice(idx.Entries, func(i,j int) bool{
    return idx.Entries[i].Path < idx.Entries[j].Path
  })

  // write the Entries
  for _,e := range idx.Entries{
    writeEntry(&buf, e)
  }

  // compute checksum of current Index 

  h := sha1.New()
  h.Write(buf.Bytes())
  buf.Write(h.Sum(nil))

  return os.WriteFile(path, buf.Bytes(), 0644)


}


// write the new entry to goon/index
func writeEntry(w io.Writer, e IndexEntry){
  binary.Write(w, binary.BigEndian, e.CTimeSec)
  binary.Write(w, binary.BigEndian, e.CTimeNSec)
    binary.Write(w, binary.BigEndian, e.MTimeSec)
    binary.Write(w, binary.BigEndian, e.MTimeNSec)
    binary.Write(w, binary.BigEndian, e.Dev)
    binary.Write(w, binary.BigEndian, e.Ino)
    binary.Write(w, binary.BigEndian, e.Mode)
    binary.Write(w, binary.BigEndian, e.UID)
    binary.Write(w, binary.BigEndian, e.GID)
    binary.Write(w, binary.BigEndian, e.FileSize)

  w.Write(e.ObjectID[:])
  binary.Write(w,binary.BigEndian, e.Flags)

  w.Write([]byte(e.Path))
  w.Write([]byte{0})

  // entry padding
  entryLen := 62 + len(e.Path) + 1
  pad := (8 - (entryLen % 8)) % 8
  w.Write(make([]byte,pad))
}
