package helper

import "unsafe"

func BytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}

	return unsafe.String(&b[0], len(b))
}

func ToFields(b []byte) []string {
	res := make([]string, 0, 16)

	i := 0

	for i < len(b) {
		for i < len(b) && isSpace(b[i]) {
			i++
		}

		start := i

		for i < len(b) && !isSpace(b[i]) {
			i++
		}

		if start < i {
			res = append(res,
				*(*string)(unsafe.Pointer(&b[start])),
			)
		}
	}

	return res
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\n' || c == '\t' || c == '\r'
}
