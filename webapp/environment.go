package webapp

import (
	"errors"
	"os"
)

// ErrMissingEnvironment is an error that represents the case when the RUNTIME environment variable is empty.
var ErrMissingEnvironment = errors.New("RUNTIME env is empty")

// Runtime struct is the parsed representation of the
// Environment in which the backend is running.
type Runtime struct {
	Environment string
}

// RuntimeFromEnv reads the value contained in the RUNTIME env variable
// and parses it in order to return a filled in Runtime struct.
func RuntimeFromEnv() (Runtime, error) {
	runtime := os.Getenv("RUNTIME")
	if len(runtime) <= 0 {
		return Runtime{}, ErrMissingEnvironment
	}

	return Runtime{Environment: runtime}, nil
}

// RuntimeFromString parses the given string in order to
// return a filled in Runtime struct.
func RuntimeFromString(runtime string) (Runtime, error) {
	if len(runtime) <= 0 {
		return Runtime{}, ErrMissingEnvironment
	}

	return Runtime{Environment: runtime}, nil
}
