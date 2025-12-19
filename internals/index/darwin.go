// go:build darwin

package index

import "syscall"

func extractMetadata(raw interface{})(cSec,cNSec,mSec,mNSec,dev,ino,uid,gid uint32){
  st := raw.(*syscall.Stat_t)
   return uint32(st.Ctimespec.Sec),
        uint32(st.Ctimespec.Nsec),
        uint32(st.Mtimespec.Sec),
        uint32(st.Mtimespec.Nsec),
        uint32(st.Dev),
        uint32(st.Ino),
        uint32(st.Uid),
        uint32(st.Gid)
}
