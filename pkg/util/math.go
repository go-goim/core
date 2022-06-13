package util

type Integer interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~int | ~int8 | ~int16 | ~int32 | ~int64
}

func Min[T Integer](a, b T) T {
	if a < b {
		return a
	}

	return b
}
