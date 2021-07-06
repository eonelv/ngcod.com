package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"

	"github.com/axgle/mahonia"

	. "ngcod.com/core"
)

const TH32CS_SNAPPROCESS = 0x00000002

type WindowsProcess struct {
	ProcessID       int
	ParentProcessID int
	Exe             string
}

func newWindowsProcess(e *syscall.ProcessEntry32) WindowsProcess {
	// Find when the string ends for decoding
	end := 0
	for {
		if e.ExeFile[end] == 0 {
			break
		}
		end++
	}
	return WindowsProcess{
		ProcessID:       int(e.ProcessID),
		ParentProcessID: int(e.ParentProcessID),
		Exe:             syscall.UTF16ToString(e.ExeFile[:end]),
	}
}

func processes() ([]WindowsProcess, error) {
	handle, err := syscall.CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer syscall.CloseHandle(handle)

	var entry syscall.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))
	// get the first process
	err = syscall.Process32First(handle, &entry)
	if err != nil {
		return nil, err
	}

	results := make([]WindowsProcess, 0, 50)
	for {
		results = append(results, newWindowsProcess(&entry))

		err = syscall.Process32Next(handle, &entry)
		if err != nil {
			// windows sends ERROR_NO_MORE_FILES on last process
			if err == syscall.ERROR_NO_MORE_FILES {
				return results, nil
			}
			return nil, err
		}
	}
}

func findProcessByName(processes []WindowsProcess, name string, parentID int) *WindowsProcess {
	for _, p := range processes {
		if bytes.Contains([]byte(strings.ToUpper(p.Exe)), []byte(strings.ToUpper(name))) {
			if parentID == 0 {
				return &p
			} else if p.ParentProcessID == parentID {
				return &p
			}
		}
	}
	return nil
}

func findProcessByID(processes []WindowsProcess, id int) *WindowsProcess {
	for _, p := range processes {
		if p.ProcessID == id {
			return &p
		}
	}
	return nil
}

func FindProcessByName(name string, parentID int) *WindowsProcess {
	prcesses, err := processes()

	if err != nil {
		LogError("Try to found out process, but failed", err)
		return nil
	}
	FindedProcess := findProcessByName(prcesses, name, parentID)

	return FindedProcess
}

func FindProcessByID(id int) *WindowsProcess {
	prcesses, err := processes()
	if err != nil {
		LogError("Try to found out process, but failed", err)
		return nil
	}
	FindedProcess := findProcessByID(prcesses, id)

	return FindedProcess
}

func FindProcessByName2(path string, name string) (bool, error) {
	cmd := exec.Command("wmic.exe", "process", "where", fmt.Sprintf(`Description='%s'`, name), "get", "ExecutablePath")

	var outerr bytes.Buffer
	cmd.Stderr = &outerr
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err == nil {
		iconv := mahonia.NewDecoder("GBK")
		data := iconv.NewReader(&out)
		result, _ := ioutil.ReadAll(data)

		outputString := string(result)
		outputString = strings.Trim(outputString, "\r\n")
		arr := strings.Split(outputString, "\r\n")
		line := ""
		for i := 0; i < len(arr); i++ {

			line = arr[i]
			line = strings.Trim(line, "\r\n")

			if line == "ExecutablePath" {
				continue
			}
			line = strings.ReplaceAll(line, "\\", "/")

			if strings.HasPrefix(strings.ToLower(line), strings.ToLower(path)) {
				return true, nil
			}
		}
	} else {
		iconv := mahonia.NewDecoder("GBK")
		data := iconv.NewReader(&outerr)
		result, _ := ioutil.ReadAll(data)
		LogDebug(string(result))
	}
	return false, nil
}
