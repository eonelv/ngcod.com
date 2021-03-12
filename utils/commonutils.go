// commonutils
package utils

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"runtime"

	"syscall"
	"unsafe"
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

func SetCmdTitle(title string) {
	systemName := runtime.GOOS
	if systemName != "windows" {
		return
	}
	//设置Title
	kernel32, _ := syscall.LoadLibrary(`kernel32.dll`)
	sct, _ := syscall.GetProcAddress(kernel32, `SetConsoleTitleW`)
	strUtf16, _ := syscall.UTF16PtrFromString(title)
	syscall.Syscall(sct, 1, uintptr(unsafe.Pointer(strUtf16)), 0, 0)
	syscall.FreeLibrary(kernel32)
}

func SetColor(color int) {
	systemName := runtime.GOOS
	if systemName != "windows" {
		return
	}
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("SetConsoleTextAttribute")
	handle, _, _ := proc.Call(uintptr(syscall.Stdout), uintptr(color))

	CloseHandle := kernel32.NewProc("CloseHandle")
	CloseHandle.Call(handle)
}

func SetCmdTitleAndColor(title string, color int) {
	SetCmdTitle(title)
	SetColor(color)
}
