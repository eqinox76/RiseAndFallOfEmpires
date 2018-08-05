package util

func Removeu32(s []uint32, i uint32) []uint32 {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
