package jack

/*
#cgo LDFLAGS: -ljack
#include <stdlib.h>
#include <jack/jack.h>

extern int goProcess(unsigned int, void *);
extern int goBufferSize(uint, void *);
extern int goSampleRate(uint, void *);
extern void goPortRegistration(jack_port_id_t, int, void *);
extern void goPortRename(jack_port_id_t, const char *, const char *, void *);
extern void goPortConnect(jack_port_id_t, jack_port_id_t, int, void *);
extern void goShutdown(void *);

jack_client_t* jack_client_open_go(const char * client_name, int options, int * status) {
	return jack_client_open(client_name, (jack_options_t) options, (jack_status_t *) status);
}

int jack_set_process_callback_go(jack_client_t * client, void * callback) {
	return jack_set_process_callback(client, goProcess, callback);
}

int jack_set_buffer_size_callback_go(jack_client_t * client, void * callback) {
	return jack_set_buffer_size_callback(client, goBufferSize, callback);
}

int jack_set_sample_rate_callback_go(jack_client_t * client, void * callback) {
	return jack_set_sample_rate_callback(client, goSampleRate, callback);
}

int jack_set_port_registration_callback_go(jack_client_t * client, void * callback) {
	return jack_set_port_registration_callback(client, goPortRegistration, callback);
}

int jack_set_port_rename_callback_go(jack_client_t * client, void * callback) {
	return jack_set_port_rename_callback(client, goPortRename, callback);
}

int jack_set_port_connect_callback_go(jack_client_t * client, void * callback) {
	return jack_set_port_connect_callback(client, goPortConnect, callback);
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
	handler                  *C.struct__jack_client
	processCallback          ProcessCallback
	bufferSizeCallback       BufferSizeCallback
	sampleRateCallback       SampleRateCallback
	portRegistrationCallback PortRegistrationCallback
	portRenameCallback       PortRenameCallback
	portConnectCallback      PortConnectCallback
	shutdownCallback         ShutdownCallback
}

type AudioSample float32

func ClientOpen(name string, options int) (*Client, int) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var status C.int
	cclient := C.jack_client_open_go(cname, C.int(options), &status)
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

func (client *Client) IsRealtime() bool {
	return C.jack_is_realtime(client.handler) != 0
}

func (client *Client) GetBufferSize() uint32 {
	return uint32(C.jack_get_buffer_size(client.handler))
}

func (client *Client) GetSampleRate() uint32 {
	return uint32(C.jack_get_sample_rate(client.handler))
}

func (client *Client) SetProcessCallback(callback ProcessCallback) int {
	client.processCallback = callback
	return int(C.jack_set_process_callback_go(client.handler, unsafe.Pointer(&client.processCallback)))
}

func (client *Client) SetBufferSizeCallback(callback BufferSizeCallback) int {
	client.bufferSizeCallback = callback
	return int(C.jack_set_buffer_size_callback_go(client.handler, unsafe.Pointer(&client.bufferSizeCallback)))
}

func (client *Client) SetSampleRateCallback(callback SampleRateCallback) int {
	client.sampleRateCallback = callback
	return int(C.jack_set_sample_rate_callback_go(client.handler, unsafe.Pointer(&client.sampleRateCallback)))
}

func (client *Client) SetPortRegistrationCallback(callback PortRegistrationCallback) int {
	client.portRegistrationCallback = callback
	return int(C.jack_set_port_registration_callback_go(client.handler, unsafe.Pointer(&client.portRegistrationCallback)))
}

func (client *Client) SetPortRenameCallback(callback PortRenameCallback) int {
	client.portRenameCallback = callback
	return int(C.jack_set_port_rename_callback_go(client.handler, unsafe.Pointer(&client.portRenameCallback)))
}

func (client *Client) SetPortConnectCallback(callback PortConnectCallback) int {
	client.portConnectCallback = callback
	return int(C.jack_set_port_connect_callback_go(client.handler, unsafe.Pointer(&client.portConnectCallback)))
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

func (client *Client) Connect(srcPort, dstPort string) int {
	csrc := C.CString(srcPort)
	defer C.free(unsafe.Pointer(csrc))
	cdst := C.CString(dstPort)
	defer C.free(unsafe.Pointer(cdst))

	return int(C.jack_connect(client.handler, csrc, cdst))
}

func (client *Client) ConnectPorts(srcPort, dstPort *Port) int {
	csrc := C.jack_port_name(srcPort.handler)
	cdst := C.jack_port_name(dstPort.handler)

	return int(C.jack_connect(client.handler, csrc, cdst))
}

func (client *Client) Disconnect(srcPort, dstPort string) int {
	csrc := C.CString(srcPort)
	defer C.free(unsafe.Pointer(csrc))
	cdst := C.CString(dstPort)
	defer C.free(unsafe.Pointer(cdst))

	return int(C.jack_disconnect(client.handler, csrc, cdst))
}

func (client *Client) DisconnectPorts(srcPort, dstPort *Port) int {
	csrc := C.jack_port_name(srcPort.handler)
	cdst := C.jack_port_name(dstPort.handler)

	return int(C.jack_disconnect(client.handler, csrc, cdst))
}

func (client *Client) GetPorts(portName, portType string, flags uint64) []string {
	cname := C.CString(portName)
	defer C.free(unsafe.Pointer(cname))
	ctype := C.CString(portType)
	defer C.free(unsafe.Pointer(ctype))

	var ports []string
	cports := C.jack_get_ports(client.handler, cname, ctype, C.ulong(flags))
	if cports != nil {
		defer C.jack_free(unsafe.Pointer(cports))
		ptr := uintptr(unsafe.Pointer(cports))
		for {
			cport := (**C.char)(unsafe.Pointer(ptr))
			if *cport == nil {
				break
			}

			str := C.GoString(*cport)
			ports = append(ports, str)
			ptr += unsafe.Sizeof(cport)
		}
	}
	return ports
}

func (client *Client) GetPortByName(name string) *Port {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	cport := C.jack_port_by_name(client.handler, cname)
	if cport != nil {
		return &Port{cport}
	}
	return nil
}

func (client *Client) GetPortById(id PortId) *Port {
	cport := C.jack_port_by_id(client.handler, C.jack_port_id_t(id))
	if cport != nil {
		return &Port{cport}
	}
	return nil
}

func (client *Client) IsPortMine(port *Port) bool {
	return C.jack_port_is_mine(client.handler, port.handler) != 0
}

type PortId uint32

type Port struct {
	handler *C.struct__jack_port
}

func (port *Port) String() string {
	return port.GetName()
}

func (port *Port) GetName() string {
	return C.GoString(C.jack_port_name(port.handler))
}

func (port *Port) GetShortName() string {
	return C.GoString(C.jack_port_short_name(port.handler))
}

func (port *Port) GetClientName() string {
	name := port.GetName()
	return name[:len(name)-len(port.GetShortName())-1]
}

func (port *Port) GetType() string {
	return C.GoString(C.jack_port_type(port.handler))
}

func (port *Port) GetBuffer(nframes uint32) []AudioSample {
	samples := C.jack_port_get_buffer(port.handler, C.jack_nframes_t(nframes))
	return (*[1 << 30]AudioSample)(samples)[:nframes:nframes]
}

func (port *Port) GetConnections() []string {
	var ports []string
	cports := C.jack_port_get_connections(port.handler)
	if cports != nil {
		defer C.jack_free(unsafe.Pointer(cports))
		ptr := uintptr(unsafe.Pointer(cports))
		for {
			cport := (**C.char)(unsafe.Pointer(ptr))
			if *cport == nil {
				break
			}

			str := C.GoString(*cport)
			ports = append(ports, str)
			ptr += unsafe.Sizeof(cport)
		}
	}
	return ports
}
