package utils

import "os"

/*
 * 文件是否想存在
 * @param path 文件目录
 */
func IsFileExists(path string) bool {
	_, err := os.Stat(path)   //用于读取文件元信息: (fileInfo, err)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}
