//Package stringutil is deprecated, plz use github.com/go-chassis/foundation/stringutil instead
package stringutil

import (
	"strings"
	"unsafe"
)

// Deprecated StringInSlice convert string to bool
// Deprecated: use github.com/go-chassis/foundation/stringutil
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Deprecated Str2bytes convert string to array of byte
// Deprecated: use github.com/go-chassis/foundation/stringutil
func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// Deprecated Bytes2str convert array of byte to string
// Deprecated: use github.com/go-chassis/foundation/stringutil
func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// Deprecated SplitToTwo split the string
// Deprecated: use github.com/go-chassis/foundation/stringutil
func SplitToTwo(s, sep string) (string, string) {
	index := strings.Index(s, sep)
	if index < 0 {
		return "", s
	}
	return s[:index], s[index+len(sep):]
}

// Deprecated SplitFirstSep split the string
// Deprecated: use github.com/go-chassis/foundation/stringutil
func SplitFirstSep(s, sep string) string {
	index := strings.Index(s, sep)
	if index < 0 {
		return ""
	}
	return s[:index]
}

// Deprecated MinInt check the minimum value of two integers
// Deprecated: use github.com/go-chassis/foundation/stringutil
func MinInt(x, y int) int {
	if x <= y {
		return x
	}

	return y
}

// Deprecated ClearStringMemory clear string memory, for very sensitive security related data
//you should clear it in memory after use
// Deprecated: use github.com/go-chassis/foundation/stringutil
func ClearStringMemory(src *string) {
	p := (*struct {
		ptr uintptr
		len int
	})(unsafe.Pointer(src))

	len := MinInt(p.len, 32)
	ptr := p.ptr
	for idx := 0; idx < len; idx = idx + 1 {
		b := (*byte)(unsafe.Pointer(&ptr))
		*b = 0
		ptr++
	}
}

// Deprecated ClearByteMemory clear byte memory, for very sensitive security related data
//you should clear it in memory after use
// Deprecated: use github.com/go-chassis/foundation/stringutil
func ClearByteMemory(src []byte) {
	len := MinInt(len(src), 32)
	for idx := 0; idx < len; idx = idx + 1 {
		src[idx] = 0
	}
}
