package jack

// #ifdef _WIN32
// #include "errno.h"
// #else
// #include <sys/errno.h>
// #endif
import "C"
import "fmt"

func StrError(status int) error {
	if 0 == status {
		return nil
	}

	var msg string
	switch status {
	case Failure:
		msg = "overall operation failed"
	case InvalidOption:
		msg = "the operation contained an invalid or unsupported option"
	case NameNotUnique:
		msg = "the desired client name was not unique"
	case ServerStarted:
		msg = "The JACK server was started as a result of this operation. Otherwise, it was running already. In either case the caller is now connected to jackd, so there is no race condition. When the server shuts down, the client will find out."
	case ServerFailed:
		msg = "unable to connect to the JACK server"
	case ServerError:
		msg = "communication error with the JACK server"
	case NoSuchClient:
		msg = "requested client does not exist"
	case LoadFailure:
		msg = "unable to load internal client"
	case InitFailure:
		msg = "unable to initialize client"
	case ShmFailure:
		msg = "unable to access shared memory"
	case VersionError:
		msg = "client's protocol version does not match"
	case BackendError:
		msg = "backend error"
	case ClientZombie:
		msg = "client zombie"
	case C.EEXIST:
		msg = "the connection is already made"
	case C.ENODATA:
		msg = "the buffer is empty"
	case C.ENOBUFS:
		msg = "there is not enough space in the buffer for the event"
	default:
		msg = fmt.Sprintf("unknown error %d", status)
	}
	return fmt.Errorf(msg)
}
