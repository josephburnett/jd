package jd

import (
	"encoding/binary"
	"hash/fnv"
)

func hash(input []byte) [8]byte {
	h := fnv.New64a()
	h.Write(input)
	var a [8]byte
	binary.LittleEndian.PutUint64(a[:], h.Sum64())
	return a
}
