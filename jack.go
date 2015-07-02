package jack

/*
#cgo LDFLAGS: -ljack
#include <stdlib.h>
#include <jack/jack.h>

extern int goProcess(unsigned int, void *);
extern void goShutdown(void *);

jack_client_t* jack_client_open_single(const char * client_name, int options, int * status) {
	return jack_client_open(client_name, (jack_options_t) options, (jack_status_t *) status);
}

int jack_set_process_callback_go(jack_client_t * client, void * callback) {
	return jack_set_process_callback(client, goProcess, callback);
}

void jack_on_shutdown_go(jack_client_t * client, void * callback) {
	jack_on_shutdown(client, goShutdown, callback);
}
*/
import "C"
import "unsafe"

const (
	// JackOptions
	NullOption    = C.JackNullOption
	NoStartServer = C.JackNoStartServer
	UseExactName  = C.JackUseExactName
	ServerName    = C.JackServerName
	LoadName      = C.JackLoadName
	LoadInit      = C.JackLoadInit
	SessionID     = C.JackSessionID

	// JackPortFlags
	PortIsInput    = C.JackPortIsInput
	PortIsOutput   = C.JackPortIsOutput
	PortIsPhysical = C.JackPortIsPhysical
	PortCanMonitor = C.JackPortCanMonitor
	PortIsTerminal = C.JackPortIsTerminal

	// JackStatus
	Failure       = C.JackFailure
	InvalidOption = C.JackInvalidOption
	NameNotUnique = C.JackNameNotUnique
	ServerStarted = C.JackServerStarted
	ServerFailed  = C.JackServerFailed
	ServerError   = C.JackServerError
	NoSuchClient  = C.JackNoSuchClient
	LoadFailure   = C.JackLoadFailure
	InitFailure   = C.JackInitFailure
	ShmFailure    = C.JackShmFailure
	VersionError  = C.JackVersionError
	BackendError  = C.JackBackendError
	ClientZombie  = C.JackClientZombie

	DEFAULT_AUDIO_TYPE = "32 bit float mono audio"
	DEFAULT_MIDI_TYPE  = "8 bit raw midi"
)

type Client struct {
	handler          *C.struct__jack_client
	processCallback  ProcessCallback
	shutdownCallback ShutdownCallback
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

func (client *Client) GetSampleRate() uint32 {
	return uint32(C.jack_get_sample_rate(client.handler))
}

func (client *Client) PortRegister(portName, portType string, flags, bufferSize uint64) *Port {
	cname := C.CString(portName)
	defer C.free(unsafe.Pointer(cname))
	ctype := C.CString(portType)
	defer C.free(unsafe.Pointer(ctype))

	cport := C.jack_port_register(client.handler, cname, ctype, C.ulong(flags), C.ulong(bufferSize))
	if cport != nil {
		return &Port{cport}
	}
	return nil
}

func (client *Client) SetProcessCallback(callback ProcessCallback) int {
	client.processCallback = callback
	return int(C.jack_set_process_callback_go(client.handler, unsafe.Pointer(&client.processCallback)))
}

func (client *Client) OnShutdown(callback ShutdownCallback) {
	client.shutdownCallback = callback
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
