package main

/*
#cgo LDFLAGS: -ljack
#include <stdlib.h>
#include <stdio.h>
#include <jack/jack.h>

extern int Process(int, void *);

jack_client_t* jack_client_open_single(const char * client_name, int options, int * status) {
	return jack_client_open(client_name, (jack_options_t) options, (jack_status_t *) status);
}

int process(uint32_t nframes, void * arg) {
	printf("%d", nframes);
	return 0;
}

int jack_set_process_callback_go(jack_client_t * client, void * callback) {
	return jack_set_process_callback(client, process, 0);
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

const (
	// JackOptions
	JackNullOption    = C.JackNullOption
	JackNoStartServer = C.JackNoStartServer
	JackUseExactName  = C.JackUseExactName
	JackServerName    = C.JackServerName
	JackLoadName      = C.JackLoadName
	JackLoadInit      = C.JackLoadInit
	JackSessionID     = C.JackSessionID

	// JackPortFlags
	JackPortIsInput    = C.JackPortIsInput
	JackPortIsOutput   = C.JackPortIsOutput
	JackPortIsPhysical = C.JackPortIsPhysical
	JackPortCanMonitor = C.JackPortCanMonitor
	JackPortIsTerminal = C.JackPortIsTerminal

	// JackStatus
	JackFailure       = C.JackFailure
	JackInvalidOption = C.JackInvalidOption
	JackNameNotUnique = C.JackNameNotUnique
	JackServerStarted = C.JackServerStarted
	JackServerFailed  = C.JackServerFailed
	JackServerError   = C.JackServerError
	JackNoSuchClient  = C.JackNoSuchClient
	JackLoadFailure   = C.JackLoadFailure
	JackInitFailure   = C.JackInitFailure
	JackShmFailure    = C.JackShmFailure
	JackVersionError  = C.JackVersionError
	JackBackendError  = C.JackBackendError
	JackClientZombie  = C.JackClientZombie

	JACK_DEFAULT_AUDIO_TYPE = "32 bit float mono audio"
	JACK_DEFAULT_MIDI_TYPE  = "8 bit raw midi"
)

type Client struct {
	handler *C.struct__jack_client
}

type Port struct {
	handler *C.struct__jack_port
}

func ClientOpen(name string, options int) (*Client, int) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var status C.int
	cclient := C.jack_client_open_single(cname, C.int(options), &status)
	var client *Client
	if cclient != nil {
		client = &Client{cclient}
	}
	return client, int(status)
}

func ClientNameSize() int {
	return int(C.jack_client_name_size())
}

func (client *Client) Activate() int {
	return int(C.jack_activate(client.handler))
}

func (client *Client) GetName() string {
	return C.GoString(C.jack_get_client_name(client.handler))
}

func (client *Client) PortRegister(portName, portType string, flags, buffer_size uint64) *Port {
	cname := C.CString(portName)
	defer C.free(unsafe.Pointer(cname))
	ctype := C.CString(portType)
	defer C.free(unsafe.Pointer(ctype))

	cport := C.jack_port_register(client.handler, cname, ctype, C.ulong(flags), C.ulong(buffer_size))
	if cport != nil {
		return &Port{cport}
	}
	return nil
}

func (client *Client) SetProcessCallback(callback func(C.int, unsafe.Pointer) C.int) int {
	return int(C.jack_set_process_callback_go(client.handler, unsafe.Pointer(&callback)))
}

func (client *Client) Close() int {
	if client == nil {
		return 0
	}
	return int(C.jack_client_close(client.handler))
}

func Process(nframes C.int, arg unsafe.Pointer) C.int {
	fmt.Println("Frames:", nframes)
	return 0
}

func main() {
	client, status := ClientOpen("Go Test", JackNoStartServer)
	if status != 0 {
		fmt.Println("Status:", status)
		return
	}
	defer client.Close()
	if code := client.SetProcessCallback(Process); code != 0 {
		fmt.Println("Failed to set process callback:", code)
		return
	}
	if code := client.Activate(); code != 0 {
		fmt.Println("Failed to activate client:", code)
		return
	}
	for port := 0; port < 2; port++ {
		fmt.Println("Registered Port:", client.PortRegister(fmt.Sprintf("port_%d", port), JACK_DEFAULT_AUDIO_TYPE, JackPortIsInput, 128))
	}
	fmt.Println(client.GetName())
	<-make(chan struct{})
}
