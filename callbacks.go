package jack

import "C"
import "unsafe"

type ProcessCallback func(uint32) int
type ShutdownCallback func()

//export goProcess
func goProcess(nframes uint, wrapper unsafe.Pointer) int {
	callback := (*ProcessCallback)(wrapper)
	return (*callback)(uint32(nframes))
}

//export goShutdown
func goShutdown(wrapper unsafe.Pointer) {
	callback := (*ShutdownCallback)(wrapper)
	(*callback)()
}
