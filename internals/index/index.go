package index

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"io"
	"os"
	"path/filepath"
	"sort"
"errors"
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
	ObjectID  [20]byte // SHA1 hash of file contents
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
	flags := uint16(len(path)) & 0x0FFF

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
		Mode:      uint32(info.Mode()),
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

	sort.Slice(idx.Entries, func(i, j int) bool {
		return idx.Entries[i].Path < idx.Entries[j].Path
	})

	// write the Entries
	for _, e := range idx.Entries {
		writeEntry(&buf, e)
	}

	// compute checksum of current Index

	h := sha1.New()
	h.Write(buf.Bytes())
	buf.Write(h.Sum(nil))

	return os.WriteFile(path, buf.Bytes(), 0o644)
}

// write the new entry to goon/index
func writeEntry(w io.Writer, e IndexEntry) {
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
	binary.Write(w, binary.BigEndian, e.Flags)

	w.Write([]byte(e.Path))
	w.Write([]byte{0})

	// entry padding
	entryLen := 62 + len(e.Path) + 1
	pad := (8 - (entryLen % 8)) % 8
	w.Write(make([]byte, pad))
}

////// ------ READING INDEX AND ENTRY

func ReadIndex(path string) (*Index, error) {
	// read ReadINdex into memory

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// create a reader to point run through the content
	r := bytes.NewReader(data)

	// instantiate index
	idx := NewIndex()

	// read index signature
	if _, err := io.ReadFull(r, idx.Signature[:]); err != nil {
		return nil, err
	}
	if string(idx.Signature[:]) != "DIRC" {
		return nil, errors.New("invalid signature")
	}

	// read version and count - next 8 bytes
	if err := binary.Read(r, binary.BigEndian, &idx.Version); err != nil {
		return nil, err
	}

	var count uint32 // 4 bytes
	if err := binary.Read(r, binary.BigEndian, &count); err != nil {
		return nil, err
	}

	// read the entries
	idx.Entries = make([]IndexEntry, 0, count)

	for i := uint32(0); i < count; i++ {
		entry, err := readEntry(r)
		if err != nil {
			return nil, err
		}

    idx.Entries = append(idx.Entries, entry)
	}

  return idx, nil
}

func readEntry(r *bytes.Reader) (IndexEntry, error) {
	var e IndexEntry

	// Fixed-size fields (62 bytes total)
	binary.Read(r, binary.BigEndian, &e.CTimeSec)
	binary.Read(r, binary.BigEndian, &e.CTimeNSec)
	binary.Read(r, binary.BigEndian, &e.MTimeSec)
	binary.Read(r, binary.BigEndian, &e.MTimeNSec)
	binary.Read(r, binary.BigEndian, &e.Dev)
	binary.Read(r, binary.BigEndian, &e.Ino)
	binary.Read(r, binary.BigEndian, &e.Mode)
	binary.Read(r, binary.BigEndian, &e.UID)
	binary.Read(r, binary.BigEndian, &e.GID)
	binary.Read(r, binary.BigEndian, &e.FileSize)

	// Object ID (20 bytes)
	r.Read(e.ObjectID[:])

	// Flags (2 bytes)
	binary.Read(r, binary.BigEndian, &e.Flags)

	// Path is a null-terminated string
	pathBytes := make([]byte, 0)
	for {
		b := make([]byte, 1)
		r.Read(b)
		if b[0] == 0 {
			break
		}
		pathBytes = append(pathBytes, b[0])
	}
	e.Path = string(pathBytes)

	// Skip padding to next 8-byte boundary
	entryLen := 62 + len(pathBytes) + 1
	pad := (8 - (entryLen % 8)) % 8
	if pad > 0 {
		r.Seek(int64(pad), io.SeekCurrent)
	}

	return e, nil
}

// remove path
func (idx *Index) RemovePath(path string){
  out := make([]IndexEntry,0, len(idx.Entries))
  for _,e := range idx.Entries{
    if e.Path != path{
      out = append(out,e)
    }
  }

  idx.Entries = out
}
