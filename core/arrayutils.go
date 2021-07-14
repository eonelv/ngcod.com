package core

import "reflect"

//将字节数组复制到对象中
//@param1 目标对象,通常也是字节数组
//@param2 源字节数组
//type MsgNetReport struct {
// ID      ObjectID
// Message [1024]byte
//}
// msgMsg := &MsgNetReport{}
//CopyArray(reflect.ValueOf(&netMsg.Message), []byte(reportMsg))
func CopyArray(dest reflect.Value, src []byte) bool {
	defer func() {
		if x := recover(); x != nil {
			LogError("CopyArray failed:", x)
		}
	}()
	return reflect.Copy(dest.Elem(), reflect.ValueOf(src)) > 0
}
