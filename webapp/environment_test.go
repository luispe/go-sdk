package webapp_test

import (
	"errors"
	"os"
	"testing"

	"github.com/pomelo-la/go-toolkit/webapp"
)

func TestRuntimeFromEnv(t *testing.T) {
	// Positive test case
	os.Setenv("RUNTIME", "production")
	defer os.Unsetenv("RUNTIME")

	runtime, err := webapp.EnvironmentFromEnvVariable()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if runtime.Name != "production" {
		t.Errorf("Unexpected environment: %s", runtime.Name)
	}

	// Negative test case
	os.Unsetenv("RUNTIME")
	_, err = webapp.EnvironmentFromEnvVariable()
	if !errors.Is(err, webapp.ErrMissingEnvironment) {
		t.Errorf("Expected error: %v, but got: %v", webapp.ErrMissingEnvironment, err)
	}
}

func TestRuntimeFromString(t *testing.T) {
	// Positive test case
	runtimeString := "exampleRuntime"
	expectedRuntime := webapp.Environment{Name: "exampleRuntime"}

	result, err := webapp.EnvironmentFromString(runtimeString)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result != expectedRuntime {
		t.Errorf("Expected runtime: %v, but got: %v", expectedRuntime, result)
	}

	// Negative test case
	emptyRuntimeString := ""
	result, err = webapp.EnvironmentFromString(emptyRuntimeString)

	if !errors.Is(webapp.ErrMissingEnvironment, err) {
		t.Errorf("Expected error: %v, but got: %v", webapp.ErrMissingEnvironment, err)
	}

	if result != (webapp.Environment{}) {
		t.Errorf("Expected empty runtime, but got: %v", result)
	}
}
