package main

import "encoding/binary"

func uint64ToBytes(i uint64) []byte {
	n1 := make([]byte, 8)
	binary.BigEndian.PutUint64(n1, i)
	return n1
}
