package lib

import (
	"errors"
	"fmt"
)

var ErrEnvironmentMisconfigured = errors.New("environment is misconfigured")

var ErrNotAValidPort = errors.New("that is not a valid port")

func EnvironmentMisconfiguredError(envVar string) error {
	return fmt.Errorf("ErrEnvironmentMisconfigured: %w : %s", ErrEnvironmentMisconfigured, envVar)
}
