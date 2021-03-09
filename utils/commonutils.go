// commonutils
package utils

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
)

func MD5(pData []byte) string {
	md5 := md5.Sum(pData)
	return string(md5[:])
}

func ReadFileData(filePath string) ([]byte, error) {
	datas, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return datas, nil
}

func CalcFileMD5(filePath string) string {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return ""
	}
	md5 := md5.Sum(bytes)
	md5Str := fmt.Sprintf("%x", md5)
	return md5Str
}
