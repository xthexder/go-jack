package main

/*
#cgo LDFLAGS: -ljack
#include <stdlib.h>
#include <jack/jack.h>

extern int Process(unsigned int, void *);
extern void Shutdown(void *);

jack_client_t* jack_client_open_single(const char * client_name, int options, int * status) {
	return jack_client_open(client_name, (jack_options_t) options, (jack_status_t *) status);
}

int jack_set_process_callback_go(jack_client_t * client, void * callback) {
	return jack_set_process_callback(client, Process, callback);
}

void jack_on_shutdown_go(jack_client_t * client, void * callback) {
	jack_on_shutdown(client, Shutdown, callback);
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
	handler          *C.struct__jack_client
	processCallback  JackProcessCallback
	shutdownCallback JackProcessCallback
}

type Port struct {
	handler *C.struct__jack_port
}

type AudioSample float32

func ClientOpen(name string, options int) (*Client, int) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var status C.int
	cclient := C.jack_client_open_single(cname, C.int(options), &status)
	var client *Client
	if cclient != nil {
		client = new(Client)
		client.handler = cclient
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

func (client *Client) SetProcessCallback(callback JackProcessCallback) int {
	client.processCallback = callback
	return int(C.jack_set_process_callback_go(client.handler, unsafe.Pointer(&client.processCallback)))
}

func (client *Client) OnShutdown(callback JackShutdownCallback) {
	C.jack_on_shutdown_go(client.handler, unsafe.Pointer(&client.shutdownCallback))
}

func (client *Client) Close() int {
	if client == nil {
		return 0
	}
	return int(C.jack_client_close(client.handler))
}

func (port *Port) GetBuffer(nframes uint32) []AudioSample {
	samples := C.jack_port_get_buffer(port.handler, C.jack_nframes_t(nframes))
	return (*[1 << 30]AudioSample)(samples)[:nframes:nframes]
}

var PortsIn []*Port
var PortsOut []*Port
var Echo []chan AudioSample

func process(nframes uint32) int {
	for i, in := range PortsIn {
		samplesIn := in.GetBuffer(nframes)
		samplesOut := PortsOut[i].GetBuffer(nframes)
		for i2, sample := range samplesIn {
			if len(Echo[i]) >= 10*1024 {
				sample += <-Echo[i] / 4
			}
			samplesOut[i2] = sample
			Echo[i] <- sample
		}
	}
	return 0
}

func shutdown() {
	fmt.Println("Shutting down")
}

func main() {
	client, status := ClientOpen("Go Test", JackNoStartServer)
	if status != 0 {
		fmt.Println("Status:", status)
		return
	}
	defer client.Close()
	if code := client.SetProcessCallback(process); code != 0 {
		fmt.Println("Failed to set process callback:", code)
		return
	}
	client.OnShutdown(shutdown)
	if code := client.Activate(); code != 0 {
		fmt.Println("Failed to activate client:", code)
		return
	}
	for i := 0; i < 2; i++ {
		portIn := client.PortRegister(fmt.Sprintf("in_%d", i), JACK_DEFAULT_AUDIO_TYPE, JackPortIsInput, 0)
		portOut := client.PortRegister(fmt.Sprintf("out_%d", i), JACK_DEFAULT_AUDIO_TYPE, JackPortIsOutput, 0)
		PortsIn = append(PortsIn, portIn)
		PortsOut = append(PortsOut, portOut)
		Echo = append(Echo, make(chan AudioSample, 10*1024))
	}
	fmt.Println(client.GetName())
	<-make(chan struct{})
}
