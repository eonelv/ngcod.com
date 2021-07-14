// multitask

package core

import (
	"runtime"
	"sync"
)

type MultiThreadTask interface {
	WriteToChannel(SrcFileDir string)
	ProcessTask(DestFileDir string)

	CreateChan()
	CloseChan()

	createWG()
	wait()
	done()
}

type BaseMultiThreadTask struct {
	wGFile *sync.WaitGroup
}

//启动MultiThreadTask. 实现该接口，让struct继承自BaseMultiThreadTask.
//实现接口函数CreateChan、CloseChan、WriteToChannel、ProcessTask
//@param1 task MultiThreadTask
//@param2 SrcFileDir 源文件目录
//@param3 DestFileDir 目标目录
func ExecTask(task MultiThreadTask, SrcFileDir string, DestFileDir string) {
	task.CreateChan()
	defer task.CloseChan()
	go task.WriteToChannel(SrcFileDir)

	task.createWG()
	for i := 0; i < runtime.NumCPU(); i++ {
		go go_ProcessTask(task, DestFileDir)
	}
	task.wait()
}

func (this *BaseMultiThreadTask) createWG() {
	this.wGFile = &sync.WaitGroup{}
	this.wGFile.Add(runtime.NumCPU())
}

func (this *BaseMultiThreadTask) wait() {
	this.wGFile.Wait()
}

func (this *BaseMultiThreadTask) done() {
	this.wGFile.Done()
}

func go_ProcessTask(task MultiThreadTask, DestFileDir string) {
	defer task.done()
	task.ProcessTask(DestFileDir)
}
