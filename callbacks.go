package jack

import "C"
import "unsafe"

type ProcessCallback func(uint32) int
type BufferSizeCallback func(uint32) int
type SampleRateCallback func(uint32) int
type PortRegistrationCallback func(PortId, bool)
type PortRenameCallback func(PortId, string, string)
type PortConnectCallback func(PortId, PortId, bool)
type ShutdownCallback func()
type ErrorFunction func(string)
type InfoFunction func(string)

//export goProcess
func goProcess(nframes uint, arg unsafe.Pointer) int {
	client := (*C.struct__jack_client)(arg)
	return clientMap[client].processCallback(uint32(nframes))
}

//export goBufferSize
func goBufferSize(nframes uint, arg unsafe.Pointer) int {
	client := (*C.struct__jack_client)(arg)
	return clientMap[client].bufferSizeCallback(uint32(nframes))
}

//export goSampleRate
func goSampleRate(nframes uint, arg unsafe.Pointer) int {
	client := (*C.struct__jack_client)(arg)
	return clientMap[client].sampleRateCallback(uint32(nframes))
}

//export goPortRegistration
func goPortRegistration(port uint, register int, arg unsafe.Pointer) {
	client := (*C.struct__jack_client)(arg)
	clientMap[client].portRegistrationCallback(PortId(port), register != 0)
}

//export goPortRename
func goPortRename(port uint, oldName, newName *C.char, arg unsafe.Pointer) {
	client := (*C.struct__jack_client)(arg)
	clientMap[client].portRenameCallback(PortId(port), C.GoString(oldName), C.GoString(newName))
}

//export goPortConnect
func goPortConnect(aport, bport uint, connect int, arg unsafe.Pointer) {
	client := (*C.struct__jack_client)(arg)
	clientMap[client].portConnectCallback(PortId(aport), PortId(bport), connect != 0)
}

//export goShutdown
func goShutdown(arg unsafe.Pointer) {
	client := (*C.struct__jack_client)(arg)
	clientMap[client].shutdownCallback()
}

//export goErrorFunction
func goErrorFunction(msg *C.char) {
	if errorFunction != nil {
		errorFunction(C.GoString(msg))
	}
}

//export goInfoFunction
func goInfoFunction(msg *C.char) {
	if infoFunction != nil {
		infoFunction(C.GoString(msg))
	}
}
