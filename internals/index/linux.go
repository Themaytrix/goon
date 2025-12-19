//go:build linux

package index

import "syscall"

func extractMetadata(raw interface{}) (cSec, cNSec, mSec, mNSec, dev, ino, uid, gid uint32) {
	st := raw.(*syscall.Stat_t)

	return uint32(st.Ctim.Sec),
		uint32(st.Ctim.Nsec),
		uint32(st.Mtim.Sec),
		uint32(st.Mtim.Nsec),
		uint32(st.Dev),
		uint32(st.Ino),
		uint32(st.Uid),
		uint32(st.Gid)
}
