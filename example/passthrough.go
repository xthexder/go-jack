package main

import (
	"fmt"

	"github.com/xthexder/go-jack"
)

var channels int = 2

var PortsIn []*jack.Port
var PortsOut []*jack.Port

func process(nframes uint32) int {
	for i, in := range PortsIn {
		samplesIn := in.GetBuffer(nframes)
		samplesOut := PortsOut[i].GetBuffer(nframes)
		for i2, sample := range samplesIn {
			samplesOut[i2] = sample
		}
	}
	return 0
}

func main() {
	client, status := jack.ClientOpen("Go Passthrough", jack.NoStartServer)
	if status != 0 {
		fmt.Println("Status:", jack.StrError(status))
		return
	}
	defer client.Close()

	for i := 0; i < channels; i++ {
		portIn := client.PortRegister(fmt.Sprintf("in_%d", i), jack.DEFAULT_AUDIO_TYPE, jack.PortIsInput, 0)
		PortsIn = append(PortsIn, portIn)
	}
	for i := 0; i < channels; i++ {
		portOut := client.PortRegister(fmt.Sprintf("out_%d", i), jack.DEFAULT_AUDIO_TYPE, jack.PortIsOutput, 0)
		PortsOut = append(PortsOut, portOut)
	}

	if code := client.SetProcessCallback(process); code != 0 {
		fmt.Println("Failed to set process callback:", jack.StrError(code))
		return
	}
	shutdown := make(chan struct{})
	client.OnShutdown(func() {
		fmt.Println("Shutting down")
		close(shutdown)
	})

	if code := client.Activate(); code != 0 {
		fmt.Println("Failed to activate client:", jack.StrError(code))
		return
	}

	fmt.Println(client.GetName())
	<-shutdown
}
