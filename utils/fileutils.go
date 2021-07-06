// fileutils
package utils

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"bytes"

	"github.com/axgle/mahonia"

	. "ngcod.com/core"
)

func PathExistAndCreate(path string) {
	if ok, _ := PathExists(path); !ok {
		os.MkdirAll(path, os.ModePerm)
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func WriteFile(data []byte, filePath string) error {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	defer f.Close()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	_, err = f.Write(data)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func CopyFile(SrcFile string, DestFile string) error {
	fileRead, err := os.Open(SrcFile)
	if err != nil {
		fmt.Println("Open err:", err)
		return err
	}
	defer fileRead.Close()

	index := strings.LastIndex(DestFile, "/")
	ParentPath := DestFile[:index]
	os.MkdirAll(ParentPath, os.ModePerm)

	//创建目标文件
	fileWrite, err := os.OpenFile(DestFile, os.O_WRONLY|os.O_CREATE, os.ModePerm)

	if err != nil {
		fmt.Println("Create err:", err)
		return err
	}
	defer fileWrite.Close()

	buf := make([]byte, 1024)
	for {
		n, err := fileRead.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := fileWrite.Write(buf[:n]); err != nil {
			return err
		}
	}
	return err
}

func Exec(cmdStr string, args ...string) error {

	var testString string = cmdStr

	for _, a := range args {
		testString += " "
		testString += a
	}
	//fmt.Println(testString)

	err := Exe_Cmd(cmdStr, true, args...)
	return err
}

func Exe_Cmd(cmdStr string, isLogInfo bool, args ...string) error {

	var testString string = cmdStr

	for _, a := range args {
		testString += " "
		testString += a
	}
	//LogDebug(testString)

	cmd := exec.Command(cmdStr, args...)
	output := make(chan []byte, 10240)
	defer close(output)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		LogDebug("get std out error:", err)
		//panic(err)
		return err
	}
	defer stdoutPipe.Close()

	stdErrorPipe, err := cmd.StderrPipe()
	if err != nil {
		LogDebug("get std err error:", err)
		//panic(err)
		return err
	}
	defer stdErrorPipe.Close()
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		iconv := mahonia.NewDecoder("GBK")
		i := 0
		for scanner.Scan() { //命令在执行的过程中, 实时地获取其输出

			var outerr bytes.Buffer
			outerr.Write(scanner.Bytes())

			readedData := iconv.NewReader(&outerr)
			result, _ := ioutil.ReadAll(readedData)
			line := string(result)

			if isLogInfo {
				LogInfo(line)
				continue
			}
			isNeedPrint := strings.Contains(line, "Error") || strings.Contains(line, "Warning")
			isNeedPrint = isNeedPrint || (i%20 == 0)
			if isNeedPrint {
				LogInfo(line)
				i = 0
			} else {
				i++
			}
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stdErrorPipe)

		iconv := mahonia.NewDecoder("GBK")

		for scanner.Scan() { // 命令在执行的过程中, 实时地获取其输出
			var outerr bytes.Buffer
			outerr.Write(scanner.Bytes())

			readedData := iconv.NewReader(&outerr)
			result, _ := ioutil.ReadAll(readedData)
			LogError(string(result))
		}
	}()

	if err := cmd.Run(); err != nil {
		LogDebug("Run ext process error:", err, cmd.Args)
		//panic(err)
		return err
	}
	return err
}
