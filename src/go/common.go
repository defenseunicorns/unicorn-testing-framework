package utf

import (
	"errors"
	"fmt"
	"os"
	"testing"

	teststructure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/require"
)

var (
	EnvVarNotFound = errors.New("expected env var not set")
)

// Setup sets up the environment where the deployment will happen.
func Setup(t *testing.T, platform interface{}) {
	teststructure.RunTestStage(t, "SETUP", func() {
		switch platform := platform.(type) {
		case EC2Platform:
			setupEC2(t, platform)
		case EKSPlatform:
			setupEKS(t, platform)
		default:
			require.Fail(t, "Unknown platform type")
		}
	})
}

// Teardown the deployment environment
func Teardown(t *testing.T, platform interface{}) {
	teststructure.RunTestStage(t, "TEARDOWN", func() {
		switch platform := platform.(type) {
		case EC2Platform:
			teardownEC2(t, platform)
		case EKSPlatform:
			teardownEKS(t, platform)
		default:
			require.Fail(t, "Unknown platform type")
		}
	})
}

// GetEnvVar gets an environment variable, returning an error if it isn't found.
func GetEnvVar(varName string) (string, error) {
	val, present := os.LookupEnv(varName)
	if !present {
		return "", fmt.Errorf("%w: %v", EnvVarNotFound, varName)
	}

	return val, nil
}
