package consistentHashing

import (
	"hash/adler32"
	"hash/crc32"
)

func GetHashFuncsMap() map[string]func(string) uint32 {
	hashFuncMaps := make(map[string]func(string) uint32)
	hashFuncMaps["hashId"] = func(id string) uint32 {
		return crc32.ChecksumIEEE([]byte(id))
	}
	hashFuncMaps["hashAdlId"] = func(id string) uint32 {
		return adler32.Checksum([]byte(id))
	}
	return hashFuncMaps
}
