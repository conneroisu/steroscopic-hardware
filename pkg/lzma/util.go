package lzma

func minUInt32(left uint32, right uint32) uint32 {
	if left < right {
		return left
	}
	return right
}
