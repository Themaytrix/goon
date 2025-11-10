package index




type Index struct{

  Signature [4]byte
  Version uint32
  Entries uint32
}

type Entry struct{
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
