package utils

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

func InsertNewLine(s string, n int) string {
	var buffer bytes.Buffer
	var n_1 = n - 1
	var l_1 = len(s) - 1
	for i, rune := range s {
		buffer.WriteRune(rune)
		if i%n == n_1 && i != l_1 {
			buffer.WriteRune('\n')
		}
	}
	return buffer.String()
}
func FileName(fileName string) (string, error) {
	AppPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	WorkPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var appConfigPath = filepath.Join(AppPath, fileName)
	if _, err = os.Stat(appConfigPath); err != nil {
		appConfigPath = filepath.Join(WorkPath, fileName)
		if _, err = os.Stat(appConfigPath); err != nil {
			return "", fmt.Errorf("配置文件不存在: %s", appConfigPath)
		}

		return appConfigPath, err
	}

	return appConfigPath, nil
}
