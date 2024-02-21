package helper

import (
	"io"
	"os"
)

// 读取文件为[]byte类型
func ReadFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()
	// 读取文件内容
	return io.ReadAll(file)
}
