# go-jack [![GoDoc](https://godoc.org/github.com/xthexder/go-jack?status.svg)](https://godoc.org/github.com/xthexder/go-jack)
Go bindings for [Jack Audio Connection Kit](http://jackaudio.org/)

## Usage

For a working passthrough example see [example/passthrough.go](https://github.com/xthexder/go-jack/blob/master/example/passthrough.go)

**Import the package:**
```go
import "github.com/xthexder/go-jack"
```

**Connect to an existing jack server:**
```go
client, _ := jack.ClientOpen("Example Client", jack.JackNoStartServer)
if client == nil {
	fmt.Println("Could not connect to jack server.")
	return
}
defer client.Close()
```

**Add a processing callback:**
```go
func process(nframes uint32) int {
	// Do processing here
	return 0
}

/* ... */

if code := client.SetProcessCallback(process); code != 0 {
	fmt.Println("Failed to set process callback.")
	return
}
```

**Activate the client:**
```go
if code := client.Activate(); code != 0 {
	fmt.Println("Failed to activate client.")
	return
}
```

**Add an output port:**
```go
port := client.PortRegister("out_1", jack.DEFAULT_AUDIO_TYPE, jack.PortIsOutput, 0)
```

**Output a sine wave:**
```go
var Port *jack.Port

func process(nframes uint32) int {
	samples := Port.GetBuffer(nframes)
	nsamples := float64(len(samples))
	for i := range samples {
		samples[i] = jack.AudioSample(math.Sin(float64(i)*math.Pi*20/nsamples) / 2)
	}
	return 0
}
```

## Implemented Bindings
 - `jack_client_t jack_client_open(client_name, options, *status)`
 - `int jack_client_close()`
 - `int jack_client_name_size()`
 - `char* jack_get_client_name(client)`
 - `void jack_on_shutdown(client, callback, arg)`
 - `int jack_set_process_callback(client, callback, arg)`
 - `jack_port_t jack_port_register(client, name, type, flags, buffer_size)`
 - `void* jack_port_get_buffer(port, nframes)`

See [Official Jack API](http://jackaudio.org/api/jack_8h.html) for detailed documentation on each of these functions.
