package webapp

import (
	"errors"
	"os"
)

// ErrMissingEnvironment is an error that represents the case when the ENVIRONMENT env variable is empty.
var ErrMissingEnvironment = errors.New("ENVIRONMENT env var is empty")

// Environment struct is the parsed representation of the
// Environment in which the backend is running.
type Environment struct {
	Name string
}

// EnvironmentFromEnvVariable reads the value contained in the ENVIRONMENT env variable
// and parses it in order to return a filled in Environment struct.
func EnvironmentFromEnvVariable() (Environment, error) {
	runtime := os.Getenv("ENVIRONMENT")
	if len(runtime) <= 0 {
		return Environment{}, ErrMissingEnvironment
	}

	return Environment{Name: runtime}, nil
}

// EnvironmentFromString parses the given string in order to
// return a filled in Environment struct.
func EnvironmentFromString(runtime string) (Environment, error) {
	if len(runtime) <= 0 {
		return Environment{}, ErrMissingEnvironment
	}

	return Environment{Name: runtime}, nil
}
