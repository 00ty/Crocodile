package utils

import (
	"fmt"
	"github.com/xluohome/phonedata"
	"strings"
	"unsafe"
)

func ByteSliceToString(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}

func StringToByteSlice(s string) []byte {
	tmp1 := (*[2]uintptr)(unsafe.Pointer(&s))
	tmp2 := [3]uintptr{tmp1[0], tmp1[1], tmp1[1]}
	return *(*[]byte)(unsafe.Pointer(&tmp2))
}

// GetPrForMobile 获取归属地
func GetPrForMobile(mobile string) (string, error) {
	pr, err := phonedata.Find(mobile)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s-%s",
		pr.Province, pr.City, pr.CardType), err
}

// SetCustomValue 判空赋值
func SetCustomValue(str string, customValue string) string {
	if str == "" {
		var builder strings.Builder
		builder.WriteString(customValue)
		return builder.String()
	}
	return str
}
