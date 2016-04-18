package main

import (
	"fmt"
	"os"

	"github.com/esqilin/go-jack"
)

var (
	portIn, portOut *jack.Port
	ch              chan string // for printing midi events
)

func process(nframes uint32) int {
	events := portIn.GetMidiEvents(nframes)
	buffer := portOut.MidiClearBuffer(nframes)
	for _, event := range events {
		ch <- fmt.Sprintf("%#v", event)
		portOut.MidiEventWrite(event, buffer)
	}

	return 0
}

func shutdown() {
	os.Exit(0)
}

func main() {
	client, status := jack.ClientOpen("Go Midi Passthrough", jack.NoStartServer)
	if status != 0 {
		fmt.Println(jack.Strerror(status))
		return
	}
	defer client.Close()

	if code := client.SetProcessCallback(process); code != 0 {
		fmt.Println("Failed to set process callback: ", jack.Strerror(code))
		return
	}
	client.OnShutdown(shutdown)

	if code := client.Activate(); code != 0 {
		fmt.Println("Failed to activate client: ", jack.Strerror(code))
		return
	}

	portIn = client.PortRegister("midi_in", jack.DEFAULT_MIDI_TYPE, jack.PortIsInput, 0)
	portOut = client.PortRegister("midi_out", jack.DEFAULT_MIDI_TYPE, jack.PortIsOutput, 0)

	fmt.Println(client.GetName())

	ch = make(chan string, 30)
	str, more := "", true
	for more {
		str, more = <-ch
		fmt.Printf("Midi Event: %s\n", str)
	}
}
