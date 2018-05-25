package jack

import "C"
import "unsafe"

type ProcessCallback func(uint32, *interface{}) int
type processCallbackWithArgs struct {
	callback ProcessCallback
	args     *interface{}
}
type BufferSizeCallback func(uint32) int
type SampleRateCallback func(uint32) int
type PortRegistrationCallback func(PortId, bool)
type PortRenameCallback func(PortId, string, string) int
type PortConnectCallback func(PortId, PortId, bool)
type ShutdownCallback func()

//export goProcess
func goProcess(nframes uint, wrapper unsafe.Pointer) int {
	callback := (*ProcessCallback)(wrapper)
	return (*callback)(uint32(nframes), nil)
}

//export goProcessWithArgs
func goProcessWithArgs(nframes uint, wrapper unsafe.Pointer) int {
	ret := (*processCallbackWithArgs)(wrapper)
	return (*ret).callback(uint32(nframes), (*ret).args)
}

//export goBufferSize
func goBufferSize(nframes uint, wrapper unsafe.Pointer) int {
	callback := (*BufferSizeCallback)(wrapper)
	return (*callback)(uint32(nframes))
}

//export goSampleRate
func goSampleRate(nframes uint, wrapper unsafe.Pointer) int {
	callback := (*SampleRateCallback)(wrapper)
	return (*callback)(uint32(nframes))
}

//export goPortRegistration
func goPortRegistration(port uint, register int, wrapper unsafe.Pointer) {
	callback := (*PortRegistrationCallback)(wrapper)
	(*callback)(PortId(port), register != 0)
}

//export goPortRename
func goPortRename(port uint, oldName, newName *C.char, wrapper unsafe.Pointer) {
	callback := (*PortRenameCallback)(wrapper)
	(*callback)(PortId(port), C.GoString(oldName), C.GoString(newName))
}

//export goPortConnect
func goPortConnect(aport, bport uint, connect int, wrapper unsafe.Pointer) {
	callback := (*PortConnectCallback)(wrapper)
	(*callback)(PortId(aport), PortId(bport), connect != 0)
}

//export goShutdown
func goShutdown(wrapper unsafe.Pointer) {
	callback := (*ShutdownCallback)(wrapper)
	(*callback)()
}
