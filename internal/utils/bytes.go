package utils

func StringToBytes32(s string) [32]byte {
	var b [32]byte
	copy(b[:], s)
	return b
}

func StringToBytes20(s string) [20]byte {
	var b [20]byte
	copy(b[:], s)
	return b
}
