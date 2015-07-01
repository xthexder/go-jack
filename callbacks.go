package main

import "C"
import "unsafe"

type JackProcessCallback func(uint32) int
type JackShutdownCallback func()

//export Process
func Process(nframes uint, wrapper unsafe.Pointer) int {
	callback := (*JackProcessCallback)(wrapper)
	return (*callback)(uint32(nframes))
}

//export Shutdown
func Shutdown(wrapper unsafe.Pointer) {
	callback := (*JackShutdownCallback)(wrapper)
	(*callback)()
}
