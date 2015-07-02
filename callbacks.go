package jack

import "C"
import "unsafe"

type JackProcessCallback func(uint32) int
type JackShutdownCallback func()

//export goProcess
func goProcess(nframes uint, wrapper unsafe.Pointer) int {
	callback := (*JackProcessCallback)(wrapper)
	return (*callback)(uint32(nframes))
}

//export goShutdown
func goShutdown(wrapper unsafe.Pointer) {
	callback := (*JackShutdownCallback)(wrapper)
	(*callback)()
}
